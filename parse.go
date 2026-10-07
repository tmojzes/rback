package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// parseRBAC parses RBAC resources from the given reader and stores them in maps under r.permissions
func (r *Rback) parseRBAC(reader io.Reader) (err error) {
	var input map[string]any

	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&input)
	if err != nil {
		return err
	}

	kind := getOrDefault[string](input, "kind", "")
	if kind != "List" {
		return fmt.Errorf("expected kind=List, but found %v", input["kind"])
	}

	r.permissions.ServiceAccounts = make(map[string]map[string]string)
	r.permissions.Roles = make(map[string]map[string]Role)
	r.permissions.RoleBindings = make(map[string]map[string]Binding)

	rawItems := getOrDefault[[]any](input, "items", nil)
	for _, i := range rawItems {
		item, ok := i.(map[string]any)
		if !ok {
			continue
		}
		nn := getNamespacedName(getMetadata(item))

		if nn.name == "" || r.shouldIgnore(nn.name) {
			continue
		}

		itemKind := getOrDefault[string](item, "kind", "")

		switch itemKind {
		case "ServiceAccount":
			if r.permissions.ServiceAccounts[nn.namespace] == nil {
				r.permissions.ServiceAccounts[nn.namespace] = make(map[string]string)
			}
			jsonStr, _ := struct2json(item)
			r.permissions.ServiceAccounts[nn.namespace][nn.name] = jsonStr
		case "RoleBinding", "ClusterRoleBinding":
			if r.permissions.RoleBindings[nn.namespace] == nil {
				r.permissions.RoleBindings[nn.namespace] = make(map[string]Binding)
			}
			r.permissions.RoleBindings[nn.namespace][nn.name] = r.toBinding(item)
		case "Role", "ClusterRole":
			if r.permissions.Roles[nn.namespace] == nil {
				r.permissions.Roles[nn.namespace] = make(map[string]Role)
			}
			r.permissions.Roles[nn.namespace][nn.name] = toRole(item)
		default:
			log.Printf("Ignoring resource kind %s", itemKind)
		}
	}
	return nil
}

func (r *Rback) shouldIgnore(name string) bool {
	for _, prefix := range r.config.ignoredPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func toKindNamespacedName(obj any) KindNamespacedName {
	o, ok := obj.(map[string]any)
	if !ok {
		return KindNamespacedName{}
	}
	return KindNamespacedName{
		kind:           getOrDefault[string](o, "kind", ""),
		NamespacedName: getNamespacedName(o),
	}
}

func getNamespacedName(metadataOrRef map[string]any) NamespacedName {
	if metadataOrRef == nil {
		return NamespacedName{}
	}
	return NamespacedName{
		namespace: getOrDefault[string](metadataOrRef, "namespace", ""),
		name:      getOrDefault[string](metadataOrRef, "name", ""),
	}
}

func getMetadata(obj map[string]any) map[string]any {
	return getOrDefault[map[string]any](obj, "metadata", nil)
}

func toRole(rawRole map[string]any) Role {
	rules := []Rule{}
	rawRules := getOrDefault[[]any](rawRole, "rules", nil)
	for _, r := range rawRules {
		rules = append(rules, toRule(r))
	}
	return Role{
		NamespacedName: getNamespacedName(getMetadata(rawRole)),
		rules:          rules,
	}
}

func (r *Rback) toBinding(rawBinding map[string]any) Binding {
	subjects := []KindNamespacedName{}
	rawSubjects := getOrDefault[[]any](rawBinding, "subjects", nil)
	for _, s := range rawSubjects {
		subject := toKindNamespacedName(s)
		if subject.name != "" && !r.shouldIgnore(subject.name) {
			subjects = append(subjects, subject)
		}
	}

	bindingNn := getNamespacedName(getMetadata(rawBinding))

	roleRef := getOrDefault[map[string]any](rawBinding, "roleRef", nil)
	role := getNamespacedName(roleRef) // note: namespace is always "", since there is no namespace field in roleRef
	if getOrDefault[string](roleRef, "kind", "") == "Role" {
		role.namespace = bindingNn.namespace
	}
	return Binding{
		NamespacedName: bindingNn,
		role:           role,
		subjects:       subjects,
	}
}

func toRule(rule any) Rule {
	r, ok := rule.(map[string]any)
	if !ok {
		return Rule{}
	}
	return Rule{
		verbs:           toStringArray(r["verbs"]),
		resources:       toStringArray(r["resources"]),
		resourceNames:   toStringArray(r["resourceNames"]),
		nonResourceURLs: toStringArray(r["nonResourceURLs"]),
		apiGroups:       toStringArray(r["apiGroups"]),
	}
}

func toStringArray(values any) []string {
	if values == nil {
		return []string{}
	}
	rawSlice, ok := values.([]any)
	if !ok {
		return []string{}
	}
	strs := make([]string, 0, len(rawSlice))
	for _, v := range rawSlice {
		if s, ok := v.(string); ok {
			strs = append(strs, s)
		}
	}
	return strs
}

// struct2json turns a map into a JSON string
func struct2json(s map[string]any) (string, error) {
	str, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(str), nil
}
