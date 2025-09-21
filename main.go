package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type NodeVisitor[T any] interface {
	VisitNode(node ast.Node) T
}

type ResultCollector[T any] interface {
	CollectResults() []T
	AddResult(item T)
}

type ItemValidator[T any] interface {
	IsValid(item T) bool
}

type ItemRenderer[T any] interface {
	RenderItem(item T) string
}

type OutputFormatter[T any] interface {
	FormatOutput(items []T) string
}

type DependencyExtractor[T any] interface {
	ExtractDependencies(item T) []string
}

type TypeNameProvider[T any] interface {
	GetTypeName(item T) string
}

type PackageProvider[T any] interface {
	GetPackage(item T) string
}

type ItemSorter[T any] interface {
	SortItems(items []T) []T
}

type DependencyResolver[T any] interface {
	ResolveDependencies(items []T) []T
}

type CodeGenerator[T any] interface {
	GenerateCode(item T) string
}

type FileWriter interface {
	WriteToFile(content string, filename string) error
}

type ImplementationNamer[T any] interface {
	GetImplementationName(item T) string
}

type LevelProvider[T any] interface {
	GetLevel(item T) int
	SetLevel(item T, level int) T
}

type GenericVisitor[T any] struct {
	nodeVisitor NodeVisitor[T]
	collector   ResultCollector[T]
	validator   ItemValidator[T]
}

func NewGenericVisitor[T any](
	nodeVisitor NodeVisitor[T],
	collector ResultCollector[T],
	validator ItemValidator[T],
) *GenericVisitor[T] {
	return &GenericVisitor[T]{
		nodeVisitor: nodeVisitor,
		collector:   collector,
		validator:   validator,
	}
}

func (gv *GenericVisitor[T]) Visit(node ast.Node) T {
	result := gv.nodeVisitor.VisitNode(node)
	if gv.validator.IsValid(result) {
		gv.collector.AddResult(result)
	}
	return result
}

func (gv *GenericVisitor[T]) GetResults() []T {
	return gv.collector.CollectResults()
}

type DependencySorter[T any] struct {
	dependencyExtractor DependencyExtractor[T]
	typeNameProvider    TypeNameProvider[T]
	dependencyResolver  DependencyResolver[T]
}

func NewDependencySorter[T any](
	extractor DependencyExtractor[T],
	nameProvider TypeNameProvider[T],
	resolver DependencyResolver[T],
) *DependencySorter[T] {
	return &DependencySorter[T]{
		dependencyExtractor: extractor,
		typeNameProvider:    nameProvider,
		dependencyResolver:  resolver,
	}
}

func (ds *DependencySorter[T]) SortItems(items []T) []T {
	return ds.dependencyResolver.ResolveDependencies(items)
}

type GenericFormatter[T any] struct {
	itemRenderer    ItemRenderer[T]
	outputFormatter OutputFormatter[T]
}

func NewGenericFormatter[T any](
	renderer ItemRenderer[T],
	formatter OutputFormatter[T],
) *GenericFormatter[T] {
	return &GenericFormatter[T]{
		itemRenderer:    renderer,
		outputFormatter: formatter,
	}
}

func (gf *GenericFormatter[T]) FormatItem(item T) string {
	return gf.itemRenderer.RenderItem(item)
}

func (gf *GenericFormatter[T]) FormatAll(items []T) string {
	return gf.outputFormatter.FormatOutput(items)
}

type GenericCodeGenerator[T any] struct {
	codeGenerator       CodeGenerator[T]
	implementationNamer ImplementationNamer[T]
	fileWriter          FileWriter
}

func NewGenericCodeGenerator[T any](
	generator CodeGenerator[T],
	namer ImplementationNamer[T],
	writer FileWriter,
) *GenericCodeGenerator[T] {
	return &GenericCodeGenerator[T]{
		codeGenerator:       generator,
		implementationNamer: namer,
		fileWriter:          writer,
	}
}

func (gcg *GenericCodeGenerator[T]) GenerateImplementation(item T) string {
	return gcg.codeGenerator.GenerateCode(item)
}

func (gcg *GenericCodeGenerator[T]) GetImplementationName(item T) string {
	return gcg.implementationNamer.GetImplementationName(item)
}

func (gcg *GenericCodeGenerator[T]) WriteToFile(content string, filename string) error {
	return gcg.fileWriter.WriteToFile(content, filename)
}

type GoStruct struct {
	Name     string
	Package  string
	Fields   []string
	Methods  []string
	Position string
	Level    int
}

type GoInterface struct {
	Name     string
	Package  string
	Methods  []string
	Position string
	Level    int
}

type GoFunction struct {
	Name       string
	Package    string
	Receiver   string
	Parameters []string
	Returns    []string
	Position   string
	Level      int
}

type StructNodeVisitor struct {
	fset *token.FileSet
	pkg  string
}

func NewStructNodeVisitor(fset *token.FileSet, pkg string) *StructNodeVisitor {
	return &StructNodeVisitor{fset: fset, pkg: pkg}
}

func (snv *StructNodeVisitor) VisitNode(node ast.Node) GoStruct {
	if ts, ok := node.(*ast.TypeSpec); ok {
		if st, ok := ts.Type.(*ast.StructType); ok {
			fields := make([]string, 0)
			if st.Fields != nil {
				for _, field := range st.Fields.List {
					if len(field.Names) > 0 {
						for _, name := range field.Names {
							fieldType := formatType(field.Type)
							fields = append(fields, fmt.Sprintf("%s %s", name.Name, fieldType))
						}
					} else {
						fieldType := formatType(field.Type)
						fields = append(fields, fieldType)
					}
				}
			}

			return GoStruct{
				Name:     ts.Name.Name,
				Package:  snv.pkg,
				Fields:   fields,
				Position: snv.fset.Position(ts.Pos()).String(),
			}
		}
	}
	return GoStruct{}
}

type StructResultCollector struct {
	results []GoStruct
}

func NewStructResultCollector() *StructResultCollector {
	return &StructResultCollector{results: make([]GoStruct, 0)}
}

func (src *StructResultCollector) CollectResults() []GoStruct {
	return src.results
}

func (src *StructResultCollector) AddResult(item GoStruct) {
	src.results = append(src.results, item)
}

type StructValidator struct{}

func (sv *StructValidator) IsValid(item GoStruct) bool {
	return item.Name != ""
}

type StructDependencyExtractor struct{}

func (sde *StructDependencyExtractor) ExtractDependencies(item GoStruct) []string {
	deps := make(map[string]bool)

	for _, field := range item.Fields {
		fieldDeps := extractTypeDependencies(field)
		for _, dep := range fieldDeps {
			if dep != item.Name {
				deps[dep] = true
			}
		}
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	return result
}

type StructTypeNameProvider struct{}

func (stnp *StructTypeNameProvider) GetTypeName(item GoStruct) string {
	return item.Name
}

type StructPackageProvider struct{}

func (spp *StructPackageProvider) GetPackage(item GoStruct) string {
	return item.Package
}

type StructItemRenderer struct{}

func (sir *StructItemRenderer) RenderItem(item GoStruct) string {
	if item.Name == "" {
		return ""
	}
	result := fmt.Sprintf("Struct: %s (Package: %s) at %s", item.Name, item.Package, item.Position)
	if len(item.Fields) > 0 {
		result += fmt.Sprintf("\n  Fields: %s", strings.Join(item.Fields, ", "))
	}
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
}

type InterfaceNodeVisitor struct {
	fset *token.FileSet
	pkg  string
}

func NewInterfaceNodeVisitor(fset *token.FileSet, pkg string) *InterfaceNodeVisitor {
	return &InterfaceNodeVisitor{fset: fset, pkg: pkg}
}

func (inv *InterfaceNodeVisitor) VisitNode(node ast.Node) GoInterface {
	if ts, ok := node.(*ast.TypeSpec); ok {
		if it, ok := ts.Type.(*ast.InterfaceType); ok {
			methods := make([]string, 0)
			if it.Methods != nil {
				for _, method := range it.Methods.List {
					if len(method.Names) > 0 {
						for _, name := range method.Names {
							if ft, ok := method.Type.(*ast.FuncType); ok {
								signature := formatFuncSignature(name.Name, ft)
								methods = append(methods, signature)
							}
						}
					} else {
						methodType := formatType(method.Type)
						methods = append(methods, methodType)
					}
				}
			}

			return GoInterface{
				Name:     ts.Name.Name,
				Package:  inv.pkg,
				Methods:  methods,
				Position: inv.fset.Position(ts.Pos()).String(),
			}
		}
	}
	return GoInterface{}
}

type InterfaceResultCollector struct {
	results []GoInterface
}

func NewInterfaceResultCollector() *InterfaceResultCollector {
	return &InterfaceResultCollector{results: make([]GoInterface, 0)}
}

func (irc *InterfaceResultCollector) CollectResults() []GoInterface {
	return irc.results
}

func (irc *InterfaceResultCollector) AddResult(item GoInterface) {
	irc.results = append(irc.results, item)
}

type InterfaceValidator struct{}

func (iv *InterfaceValidator) IsValid(item GoInterface) bool {
	return item.Name != ""
}

type InterfaceDependencyExtractor struct{}

func (ide *InterfaceDependencyExtractor) ExtractDependencies(item GoInterface) []string {
	deps := make(map[string]bool)

	for _, method := range item.Methods {
		methodDeps := extractTypeDependencies(method)
		for _, dep := range methodDeps {
			if dep != item.Name {
				deps[dep] = true
			}
		}
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	return result
}

type InterfaceTypeNameProvider struct{}

func (itnp *InterfaceTypeNameProvider) GetTypeName(item GoInterface) string {
	return item.Name
}

type InterfacePackageProvider struct{}

func (ipp *InterfacePackageProvider) GetPackage(item GoInterface) string {
	return item.Package
}

type InterfaceItemRenderer struct{}

func (iir *InterfaceItemRenderer) RenderItem(item GoInterface) string {
	if item.Name == "" {
		return ""
	}
	result := fmt.Sprintf("Interface: %s (Package: %s) at %s", item.Name, item.Package, item.Position)
	if len(item.Methods) > 0 {
		result += fmt.Sprintf("\n  Methods: %s", strings.Join(item.Methods, ", "))
	}
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
}

type InterfaceNoOpCodeGenerator struct{}

func (incg *InterfaceNoOpCodeGenerator) GenerateCode(item GoInterface) string {
	if item.Name == "" {
		return ""
	}

	implName := fmt.Sprintf("NoOp%s", item.Name)
	var builder strings.Builder

	// Add comment header
	builder.WriteString(fmt.Sprintf("// %s is a no-op implementation of %s interface (Level %d)\n", implName, item.Name, item.Level))

	// Struct definition with level field
	builder.WriteString(fmt.Sprintf("type %s struct {\n", implName))
	builder.WriteString(fmt.Sprintf("\tlevel int // Dependency level: %d\n", item.Level))
	builder.WriteString("}\n\n")

	// Constructor
	builder.WriteString(fmt.Sprintf("// New%s creates a new no-op implementation at the specified level\n", implName))
	builder.WriteString(fmt.Sprintf("func New%s(level int) *%s {\n", implName, implName))
	builder.WriteString(fmt.Sprintf("\treturn &%s{level: level}\n", implName))
	builder.WriteString("}\n\n")

	// GetLevel method
	builder.WriteString(fmt.Sprintf("// GetLevel returns the dependency level of this %s\n", implName))
	builder.WriteString(fmt.Sprintf("func (n *%s) GetLevel() int {\n", implName))
	builder.WriteString("\treturn n.level\n")
	builder.WriteString("}\n\n")

	// Generate methods
	for _, method := range item.Methods {
		methodImpl := generateMethodImplementation(method, implName, item.Level)
		if methodImpl != "" {
			builder.WriteString(methodImpl)
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

type InterfaceImplementationNamer struct{}

func (iin *InterfaceImplementationNamer) GetImplementationName(item GoInterface) string {
	return fmt.Sprintf("NoOp%s", item.Name)
}

type TopologicalDependencyResolver[T any] struct {
	dependencyExtractor DependencyExtractor[T]
	typeNameProvider    TypeNameProvider[T]
}

func NewTopologicalDependencyResolver[T any](
	extractor DependencyExtractor[T],
	nameProvider TypeNameProvider[T],
) *TopologicalDependencyResolver[T] {
	return &TopologicalDependencyResolver[T]{
		dependencyExtractor: extractor,
		typeNameProvider:    nameProvider,
	}
}

func (tdr *TopologicalDependencyResolver[T]) ResolveDependencies(items []T) []T {
	// Build the item map
	allItems := make(map[string]T)
	for _, item := range items {
		name := tdr.typeNameProvider.GetTypeName(item)
		if name != "" {
			allItems[name] = item
		}
	}

	// Build dependency graph
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all nodes
	for name := range allItems {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// Build edges (dependencies)
	for name, item := range allItems {
		deps := tdr.dependencyExtractor.ExtractDependencies(item)
		for _, dep := range deps {
			if _, exists := allItems[dep]; exists {
				graph[dep] = append(graph[dep], name)
				inDegree[name]++
			}
		}
	}

	// Kahn's algorithm for topological sorting
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	result := make([]T, 0, len(items))
	processed := make(map[string]bool)

	for len(queue) > 0 {
		sort.Strings(queue)

		current := queue[0]
		queue = queue[1:]

		if processed[current] {
			continue
		}
		processed[current] = true

		if item, exists := allItems[current]; exists {
			result = append(result, item)
		}

		for _, dependent := range graph[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// Add any remaining items
	for _, item := range items {
		name := tdr.typeNameProvider.GetTypeName(item)
		if !processed[name] {
			result = append(result, item)
		}
	}

	return result
}

type SimpleOutputFormatter[T any] struct{}

func (sof *SimpleOutputFormatter[T]) FormatOutput(items []T) string {
	return "" // Individual items are formatted by ItemRenderer
}

type SimpleFileWriter struct{}

func (sfw *SimpleFileWriter) WriteToFile(content string, filename string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

type AnalysisEngine[T any] struct {
	visitor       *GenericVisitor[T]
	sorter        ItemSorter[T]
	formatter     *GenericFormatter[T]
	codeGenerator *GenericCodeGenerator[T]
}

func NewAnalysisEngine[T any](
	visitor *GenericVisitor[T],
	sorter ItemSorter[T],
	formatter *GenericFormatter[T],
	codeGenerator *GenericCodeGenerator[T],
) *AnalysisEngine[T] {
	return &AnalysisEngine[T]{
		visitor:       visitor,
		sorter:        sorter,
		formatter:     formatter,
		codeGenerator: codeGenerator,
	}
}

func (ae *AnalysisEngine[T]) Analyze(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		ae.visitor.Visit(n)
		return true
	})
}

func (ae *AnalysisEngine[T]) GetSortedResults() []T {
	results := ae.visitor.GetResults()
	if ae.sorter != nil {
		results = ae.sorter.SortItems(results)
	}
	return results
}

func (ae *AnalysisEngine[T]) PrintResults() {
	results := ae.GetSortedResults()

	for level, result := range results {
		formatted := ae.formatter.FormatItem(result)
		if formatted != "" {
			fmt.Printf("[Level %d] %s\n", level, formatted)

			// Generate NoOp if available
			if ae.codeGenerator != nil {
				if noopImpl := ae.codeGenerator.GenerateImplementation(result); noopImpl != "" {
					fmt.Printf("\n--- NoOp Implementation ---\n")
					fmt.Println(noopImpl)
				}
			}
		}
	}
}

func (ae *AnalysisEngine[T]) GenerateCodeFile(filename string) error {
	if ae.codeGenerator == nil {
		return fmt.Errorf("code generator not available")
	}

	results := ae.GetSortedResults()

	var builder strings.Builder
	builder.WriteString("// Code generated by go-ast-analyzer; DO NOT EDIT.\n\n")
	builder.WriteString("package main\n\n")

	for _, result := range results {
		if code := ae.codeGenerator.GenerateImplementation(result); code != "" {
			builder.WriteString(code)
			builder.WriteString("\n")
		}
	}

	return ae.codeGenerator.WriteToFile(builder.String(), filename)
}

func extractTypeDependencies(typeStr string) []string {
	deps := make(map[string]bool)

	cleaned := strings.ReplaceAll(typeStr, "*", "")
	cleaned = strings.ReplaceAll(cleaned, "[]", "")
	cleaned = strings.ReplaceAll(cleaned, "map[", "")
	cleaned = strings.ReplaceAll(cleaned, "chan ", "")
	cleaned = strings.ReplaceAll(cleaned, "<-", "")

	words := strings.FieldsFunc(cleaned, func(c rune) bool {
		return c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}' ||
			c == ',' || c == ' ' || c == '\t' || c == '\n'
	})

	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		if isBuiltinType(word) {
			continue
		}

		if strings.Contains(word, ".") {
			parts := strings.Split(word, ".")
			if len(parts) == 2 {
				deps[parts[1]] = true
			}
		} else if isValidIdentifier(word) {
			deps[word] = true
		}
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}

	return result
}

func isBuiltinType(name string) bool {
	builtins := map[string]bool{
		"bool": true, "byte": true, "complex64": true, "complex128": true,
		"error": true, "float32": true, "float64": true, "int": true,
		"int8": true, "int16": true, "int32": true, "int64": true,
		"rune": true, "string": true, "uint": true, "uint8": true,
		"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
		"interface": true, "func": true, "struct": true,
	}
	return builtins[name]
}

func isValidIdentifier(name string) bool {
	if name == "" {
		return false
	}

	first := rune(name[0])
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}

	for _, r := range name[1:] {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	return true
}

func formatType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + formatType(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + formatType(t.Elt)
		}
		return fmt.Sprintf("[%s]%s", formatExpr(t.Len), formatType(t.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", formatType(t.Key), formatType(t.Value))
	case *ast.FuncType:
		return formatFuncType(t)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", formatType(t.X), t.Sel.Name)
	default:
		return "unknown"
	}
}

func formatFuncType(ft *ast.FuncType) string {
	params := ""
	if ft.Params != nil {
		paramList := make([]string, 0)
		for _, param := range ft.Params.List {
			paramType := formatType(param.Type)
			if len(param.Names) > 0 {
				for range param.Names {
					paramList = append(paramList, paramType)
				}
			} else {
				paramList = append(paramList, paramType)
			}
		}
		params = strings.Join(paramList, ", ")
	}

	results := ""
	if ft.Results != nil {
		resultList := make([]string, 0)
		for _, result := range ft.Results.List {
			resultType := formatType(result.Type)
			if len(result.Names) > 0 {
				for range result.Names {
					resultList = append(resultList, resultType)
				}
			} else {
				resultList = append(resultList, resultType)
			}
		}
		if len(resultList) == 1 {
			results = " " + resultList[0]
		} else if len(resultList) > 1 {
			results = " (" + strings.Join(resultList, ", ") + ")"
		}
	}

	return fmt.Sprintf("func(%s)%s", params, results)
}

func formatFuncSignature(name string, ft *ast.FuncType) string {
	return name + strings.TrimPrefix(formatFuncType(ft), "func")
}

func formatExpr(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Value
	case *ast.Ident:
		return e.Name
	default:
		return "complex_expr"
	}
}

func generateMethodImplementation(methodSig, implName string, level int) string {
	methodName, params, returns := parseMethodSignature(methodSig)
	if methodName == "" {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("// %s is a no-op implementation (Level %d)\n", methodName, level))
	builder.WriteString(fmt.Sprintf("func (n *%s) %s(%s)", implName, methodName, params))

	if returns != "" {
		builder.WriteString(fmt.Sprintf(" %s", returns))
	}

	builder.WriteString(" {\n")
	builder.WriteString(fmt.Sprintf("\t// TODO: Implement %s (Level %d)\n", methodName, level))

	if returns != "" {
		zeroValues := generateZeroValues(returns)
		if zeroValues != "" {
			builder.WriteString(fmt.Sprintf("\treturn %s\n", zeroValues))
		}
	}

	builder.WriteString("}")
	return builder.String()
}

func parseMethodSignature(methodSig string) (name, params, returns string) {
	parts := strings.Split(methodSig, "(")
	if len(parts) < 2 {
		return "", "", ""
	}

	name = strings.TrimSpace(parts[0])
	remainder := strings.Join(parts[1:], "(")

	parenCount := 0
	paramEnd := -1

	for i, char := range remainder {
		if char == '(' {
			parenCount++
		} else if char == ')' {
			if parenCount == 0 {
				paramEnd = i
				break
			}
			parenCount--
		}
	}

	if paramEnd == -1 {
		return name, "", ""
	}

	params = remainder[:paramEnd]
	returns = strings.TrimSpace(remainder[paramEnd+1:])

	return name, params, returns
}

func generateZeroValues(returnTypes string) string {
	if returnTypes == "" {
		return ""
	}

	returnTypes = strings.Trim(returnTypes, "()")
	types := strings.Split(returnTypes, ",")
	zeroVals := make([]string, 0, len(types))

	for _, typ := range types {
		typ = strings.TrimSpace(typ)
		if strings.Contains(typ, " ") {
			parts := strings.Fields(typ)
			if len(parts) >= 2 {
				typ = parts[len(parts)-1]
			}
		}

		zeroVal := getZeroValue(typ)
		zeroVals = append(zeroVals, zeroVal)
	}

	return strings.Join(zeroVals, ", ")
}

func getZeroValue(typ string) string {
	if strings.HasPrefix(typ, "*") || strings.HasPrefix(typ, "[]") || strings.HasPrefix(typ, "map[") || strings.Contains(typ, "chan") {
		return "nil"
	}

	switch typ {
	case "bool":
		return "false"
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "complex64", "complex128":
		return "0+0i"
	case "error":
		return "nil"
	default:
		if strings.Contains(typ, ".") {
			return "nil"
		}
		return fmt.Sprintf("%s{}", typ)
	}
}

func main() {
	var (
		dirs        = flag.String("dirs", ".", "Comma-separated list of directories to analyze")
		showStructs = flag.Bool("structs", false, "Show structs")
		showIfaces  = flag.Bool("interfaces", false, "Show interfaces")
		genNoOp     = flag.Bool("noop", false, "Generate NoOp implementations for interfaces")
		noOpDir     = flag.String("noop-dir", "./noop", "Directory to save NoOp implementations")
	)

	flag.Parse()

	if *genNoOp && *noOpDir != "" {
		if err := os.MkdirAll(*noOpDir, 0755); err != nil {
			log.Fatalf("Failed to create NoOp directory %s: %v", *noOpDir, err)
		}
	}

	directories := strings.Split(*dirs, ",")

	for _, dir := range directories {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}

		if err := walkDirectory(dir, *showStructs, *showIfaces, *genNoOp, *noOpDir); err != nil {
			log.Printf("Error analyzing directory %s: %v", dir, err)
		}
	}
}

func walkDirectory(dir string, showStructs, showIfaces, genNoOp bool, noOpDir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			return processFile(path, showStructs, showIfaces, genNoOp, noOpDir)
		}

		return nil
	})
}

func processFile(filename string, showStructs, showIfaces, genNoOp bool, noOpDir string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %v", filename, err)
	}

	pkg := node.Name.Name
	fmt.Printf("\n=== Analyzing file: %s ===\n", filename)

	if showStructs {
		fmt.Println("\n--- Structs (Dependency Order) ---")

		// Build struct analysis engine with segregated interfaces
		structVisitor := NewGenericVisitor(
			NewStructNodeVisitor(fset, pkg),
			NewStructResultCollector(),
			&StructValidator{},
		)

		structSorter := NewDependencySorter(
			&StructDependencyExtractor{},
			&StructTypeNameProvider{},
			NewTopologicalDependencyResolver(
				&StructDependencyExtractor{},
				&StructTypeNameProvider{},
			),
		)

		structFormatter := NewGenericFormatter(
			&StructItemRenderer{},
			&SimpleOutputFormatter[GoStruct]{},
		)

		structEngine := NewAnalysisEngine(
			structVisitor,
			structSorter,
			structFormatter,
			nil, // No code generator for structs
		)

		structEngine.Analyze(node)
		structEngine.PrintResults()
	}

	if showIfaces {
		fmt.Println("\n--- Interfaces (Dependency Order) ---")

		// Build interface analysis engine with segregated interfaces
		interfaceVisitor := NewGenericVisitor(
			NewInterfaceNodeVisitor(fset, pkg),
			NewInterfaceResultCollector(),
			&InterfaceValidator{},
		)

		interfaceSorter := NewDependencySorter(
			&InterfaceDependencyExtractor{},
			&InterfaceTypeNameProvider{},
			NewTopologicalDependencyResolver(
				&InterfaceDependencyExtractor{},
				&InterfaceTypeNameProvider{},
			),
		)

		interfaceFormatter := NewGenericFormatter(
			&InterfaceItemRenderer{},
			&SimpleOutputFormatter[GoInterface]{},
		)

		var interfaceCodeGen *GenericCodeGenerator[GoInterface]
		if genNoOp {
			interfaceCodeGen = NewGenericCodeGenerator(
				&InterfaceNoOpCodeGenerator{},
				&InterfaceImplementationNamer{},
				&SimpleFileWriter{},
			)
		}

		interfaceEngine := NewAnalysisEngine(
			interfaceVisitor,
			interfaceSorter,
			interfaceFormatter,
			interfaceCodeGen,
		)

		interfaceEngine.Analyze(node)
		interfaceEngine.PrintResults()

		// Generate NoOp file if requested
		if genNoOp && noOpDir != "" {
			baseFilename := filepath.Base(filename)
			noOpFilename := filepath.Join(noOpDir, "noop_"+strings.TrimSuffix(baseFilename, ".go")+"_interfaces.go")
			if err := interfaceEngine.GenerateCodeFile(noOpFilename); err != nil {
				log.Printf("Failed to generate NoOp file %s: %v", noOpFilename, err)
			} else {
				fmt.Printf("Generated NoOp implementations: %s\n", noOpFilename)
			}
		}
	}

	return nil
}
