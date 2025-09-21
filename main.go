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

type GoVariable struct {
	Name     string
	Package  string
	Type     string
	Position string
	Level    int
}

type GoConstant struct {
	Name     string
	Package  string
	Type     string
	Value    string
	Position string
	Level    int
}

type GoImport struct {
	Name     string
	Path     string
	Position string
	Level    int
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

type FunctionNodeVisitor struct {
	fset *token.FileSet
	pkg  string
}

func NewFunctionNodeVisitor(fset *token.FileSet, pkg string) *FunctionNodeVisitor {
	return &FunctionNodeVisitor{fset: fset, pkg: pkg}
}

func (fnv *FunctionNodeVisitor) VisitNode(node ast.Node) GoFunction {
	if fn, ok := node.(*ast.FuncDecl); ok {
		receiver := ""
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			receiver = formatType(fn.Recv.List[0].Type)
		}

		params := make([]string, 0)
		if fn.Type.Params != nil {
			for _, param := range fn.Type.Params.List {
				paramType := formatType(param.Type)
				if len(param.Names) > 0 {
					for _, name := range param.Names {
						params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
					}
				} else {
					params = append(params, paramType)
				}
			}
		}

		returns := make([]string, 0)
		if fn.Type.Results != nil {
			for _, result := range fn.Type.Results.List {
				resultType := formatType(result.Type)
				if len(result.Names) > 0 {
					for _, name := range result.Names {
						returns = append(returns, fmt.Sprintf("%s %s", name.Name, resultType))
					}
				} else {
					returns = append(returns, resultType)
				}
			}
		}

		return GoFunction{
			Name:       fn.Name.Name,
			Package:    fnv.pkg,
			Receiver:   receiver,
			Parameters: params,
			Returns:    returns,
			Position:   fnv.fset.Position(fn.Pos()).String(),
		}
	}
	return GoFunction{}
}

type FunctionResultCollector struct {
	results []GoFunction
}

func NewFunctionResultCollector() *FunctionResultCollector {
	return &FunctionResultCollector{results: make([]GoFunction, 0)}
}

func (frc *FunctionResultCollector) CollectResults() []GoFunction {
	return frc.results
}

func (frc *FunctionResultCollector) AddResult(item GoFunction) {
	frc.results = append(frc.results, item)
}

type FunctionValidator struct{}

func (fv *FunctionValidator) IsValid(item GoFunction) bool {
	return item.Name != ""
}

type FunctionDependencyExtractor struct{}

func (fde *FunctionDependencyExtractor) ExtractDependencies(item GoFunction) []string {
	deps := make(map[string]bool)

	// Dependencies from receiver
	if item.Receiver != "" {
		receiverDeps := extractTypeDependencies(item.Receiver)
		for _, dep := range receiverDeps {
			deps[dep] = true
		}
	}

	// Dependencies from parameters
	for _, param := range item.Parameters {
		paramDeps := extractTypeDependencies(param)
		for _, dep := range paramDeps {
			deps[dep] = true
		}
	}

	// Dependencies from return types
	for _, ret := range item.Returns {
		retDeps := extractTypeDependencies(ret)
		for _, dep := range retDeps {
			deps[dep] = true
		}
	}

	result := make([]string, 0, len(deps))
	for dep := range deps {
		result = append(result, dep)
	}
	return result
}

type FunctionTypeNameProvider struct{}

func (ftnp *FunctionTypeNameProvider) GetTypeName(item GoFunction) string {
	if item.Receiver != "" {
		return fmt.Sprintf("%s.%s", item.Receiver, item.Name)
	}
	return item.Name
}

type FunctionPackageProvider struct{}

func (fpp *FunctionPackageProvider) GetPackage(item GoFunction) string {
	return item.Package
}

type FunctionItemRenderer struct{}

func (fir *FunctionItemRenderer) RenderItem(item GoFunction) string {
	if item.Name == "" {
		return ""
	}
	result := fmt.Sprintf("Function: %s (Package: %s) at %s", item.Name, item.Package, item.Position)
	if item.Receiver != "" {
		result = fmt.Sprintf("Method: %s (Receiver: %s, Package: %s) at %s", item.Name, item.Receiver, item.Package, item.Position)
	}
	if len(item.Parameters) > 0 {
		result += fmt.Sprintf("\n  Parameters: %s", strings.Join(item.Parameters, ", "))
	}
	if len(item.Returns) > 0 {
		result += fmt.Sprintf("\n  Returns: %s", strings.Join(item.Returns, ", "))
	}
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
}

type VariableNodeVisitor struct {
	fset *token.FileSet
	pkg  string
}

func NewVariableNodeVisitor(fset *token.FileSet, pkg string) *VariableNodeVisitor {
	return &VariableNodeVisitor{fset: fset, pkg: pkg}
}

func (vnv *VariableNodeVisitor) VisitNode(node ast.Node) GoVariable {
	if vs, ok := node.(*ast.ValueSpec); ok {
		for i, name := range vs.Names {
			varType := ""
			if vs.Type != nil {
				varType = formatType(vs.Type)
			} else if len(vs.Values) > i {
				varType = "inferred"
			}

			return GoVariable{
				Name:     name.Name,
				Package:  vnv.pkg,
				Type:     varType,
				Position: vnv.fset.Position(name.Pos()).String(),
			}
		}
	}
	return GoVariable{}
}

type VariableResultCollector struct {
	results []GoVariable
}

func NewVariableResultCollector() *VariableResultCollector {
	return &VariableResultCollector{results: make([]GoVariable, 0)}
}

func (vrc *VariableResultCollector) CollectResults() []GoVariable {
	return vrc.results
}

func (vrc *VariableResultCollector) AddResult(item GoVariable) {
	vrc.results = append(vrc.results, item)
}

type VariableValidator struct{}

func (vv *VariableValidator) IsValid(item GoVariable) bool {
	return item.Name != ""
}

type VariableDependencyExtractor struct{}

func (vde *VariableDependencyExtractor) ExtractDependencies(item GoVariable) []string {
	return extractTypeDependencies(item.Type)
}

type VariableTypeNameProvider struct{}

func (vtnp *VariableTypeNameProvider) GetTypeName(item GoVariable) string {
	return item.Name
}

type VariablePackageProvider struct{}

func (vpp *VariablePackageProvider) GetPackage(item GoVariable) string {
	return item.Package
}

type VariableItemRenderer struct{}

func (vir *VariableItemRenderer) RenderItem(item GoVariable) string {
	if item.Name == "" {
		return ""
	}
	result := fmt.Sprintf("Variable: %s %s (Package: %s) at %s", item.Name, item.Type, item.Package, item.Position)
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
}

type ConstantNodeVisitor struct {
	fset *token.FileSet
	pkg  string
}

func NewConstantNodeVisitor(fset *token.FileSet, pkg string) *ConstantNodeVisitor {
	return &ConstantNodeVisitor{fset: fset, pkg: pkg}
}

func (cnv *ConstantNodeVisitor) VisitNode(node ast.Node) GoConstant {
	if vs, ok := node.(*ast.ValueSpec); ok {
		for i, name := range vs.Names {
			constType := ""
			if vs.Type != nil {
				constType = formatType(vs.Type)
			}

			value := ""
			if len(vs.Values) > i {
				value = formatExpr(vs.Values[i])
			}

			return GoConstant{
				Name:     name.Name,
				Package:  cnv.pkg,
				Type:     constType,
				Value:    value,
				Position: cnv.fset.Position(name.Pos()).String(),
			}
		}
	}
	return GoConstant{}
}

type ConstantResultCollector struct {
	results []GoConstant
}

func NewConstantResultCollector() *ConstantResultCollector {
	return &ConstantResultCollector{results: make([]GoConstant, 0)}
}

func (crc *ConstantResultCollector) CollectResults() []GoConstant {
	return crc.results
}

func (crc *ConstantResultCollector) AddResult(item GoConstant) {
	crc.results = append(crc.results, item)
}

type ConstantValidator struct{}

func (cv *ConstantValidator) IsValid(item GoConstant) bool {
	return item.Name != ""
}

type ConstantDependencyExtractor struct{}

func (cde *ConstantDependencyExtractor) ExtractDependencies(item GoConstant) []string {
	return extractTypeDependencies(item.Type)
}

type ConstantTypeNameProvider struct{}

func (ctnp *ConstantTypeNameProvider) GetTypeName(item GoConstant) string {
	return item.Name
}

type ConstantPackageProvider struct{}

func (cpp *ConstantPackageProvider) GetPackage(item GoConstant) string {
	return item.Package
}

type ConstantItemRenderer struct{}

func (cir *ConstantItemRenderer) RenderItem(item GoConstant) string {
	if item.Name == "" {
		return ""
	}
	result := fmt.Sprintf("Constant: %s", item.Name)
	if item.Type != "" {
		result += fmt.Sprintf(" %s", item.Type)
	}
	if item.Value != "" {
		result += fmt.Sprintf(" = %s", item.Value)
	}
	result += fmt.Sprintf(" (Package: %s) at %s", item.Package, item.Position)
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
}

type ImportNodeVisitor struct {
	fset *token.FileSet
}

func NewImportNodeVisitor(fset *token.FileSet) *ImportNodeVisitor {
	return &ImportNodeVisitor{fset: fset}
}

func (inv *ImportNodeVisitor) VisitNode(node ast.Node) GoImport {
	if is, ok := node.(*ast.ImportSpec); ok {
		name := ""
		if is.Name != nil {
			name = is.Name.Name
		}

		path := ""
		if is.Path != nil {
			path = is.Path.Value
		}

		return GoImport{
			Name:     name,
			Path:     path,
			Position: inv.fset.Position(is.Pos()).String(),
		}
	}
	return GoImport{}
}

type ImportResultCollector struct {
	results []GoImport
}

func NewImportResultCollector() *ImportResultCollector {
	return &ImportResultCollector{results: make([]GoImport, 0)}
}

func (irc *ImportResultCollector) CollectResults() []GoImport {
	return irc.results
}

func (irc *ImportResultCollector) AddResult(item GoImport) {
	irc.results = append(irc.results, item)
}

type ImportValidator struct{}

func (iv *ImportValidator) IsValid(item GoImport) bool {
	return item.Path != ""
}

type ImportDependencyExtractor struct{}

func (ide *ImportDependencyExtractor) ExtractDependencies(item GoImport) []string {
	return []string{} // Imports don't have dependencies in our model
}

type ImportTypeNameProvider struct{}

func (itnp *ImportTypeNameProvider) GetTypeName(item GoImport) string {
	return item.Path
}

type ImportPackageProvider struct{}

func (ipp *ImportPackageProvider) GetPackage(item GoImport) string {
	return item.Path
}

type ImportItemRenderer struct{}

func (iir *ImportItemRenderer) RenderItem(item GoImport) string {
	if item.Path == "" {
		return ""
	}
	result := fmt.Sprintf("Import: %s", item.Path)
	if item.Name != "" && item.Name != "." {
		result = fmt.Sprintf("Import: %s as %s", item.Path, item.Name)
	}
	result += fmt.Sprintf(" at %s", item.Position)
	if item.Level > 0 {
		result += fmt.Sprintf("\n  Level: %d", item.Level)
	}
	return result
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

type AlphabeticalDependencyResolver[T any] struct {
	typeNameProvider TypeNameProvider[T]
}

func NewAlphabeticalDependencyResolver[T any](
	nameProvider TypeNameProvider[T],
) *AlphabeticalDependencyResolver[T] {
	return &AlphabeticalDependencyResolver[T]{
		typeNameProvider: nameProvider,
	}
}

func (adr *AlphabeticalDependencyResolver[T]) ResolveDependencies(items []T) []T {
	result := make([]T, len(items))
	copy(result, items)

	sort.Slice(result, func(i, j int) bool {
		nameI := adr.typeNameProvider.GetTypeName(result[i])
		nameJ := adr.typeNameProvider.GetTypeName(result[j])
		return nameI < nameJ
	})

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
		"interface": true, "func": true, "struct": true, "any": true,
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
	case *ast.ChanType:
		dir := ""
		if t.Dir == ast.SEND {
			dir = "<-"
		} else if t.Dir == ast.RECV {
			dir = "<-"
		}
		return fmt.Sprintf("%schan %s", dir, formatType(t.Value))
	case *ast.FuncType:
		return formatFuncType(t)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", formatType(t.X), t.Sel.Name)
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", formatType(t.X), formatType(t.Index))
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
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", formatExpr(e.X), e.Op.String(), formatExpr(e.Y))
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", e.Op.String(), formatExpr(e.X))
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

func analyzeDecl(decl ast.Decl, engines map[string]interface{}) {
	switch d := decl.(type) {
	case *ast.GenDecl:
		switch d.Tok {
		case token.TYPE:
			for _, spec := range d.Specs {
				if engine, ok := engines["structs"].(*AnalysisEngine[GoStruct]); ok {
					engine.Analyze(spec)
				}
				if engine, ok := engines["interfaces"].(*AnalysisEngine[GoInterface]); ok {
					engine.Analyze(spec)
				}
			}
		case token.VAR:
			if engine, ok := engines["variables"].(*AnalysisEngine[GoVariable]); ok {
				for _, spec := range d.Specs {
					engine.Analyze(spec)
				}
			}
		case token.CONST:
			if engine, ok := engines["constants"].(*AnalysisEngine[GoConstant]); ok {
				for _, spec := range d.Specs {
					engine.Analyze(spec)
				}
			}
		case token.IMPORT:
			if engine, ok := engines["imports"].(*AnalysisEngine[GoImport]); ok {
				for _, spec := range d.Specs {
					engine.Analyze(spec)
				}
			}
		}
	case *ast.FuncDecl:
		if engine, ok := engines["functions"].(*AnalysisEngine[GoFunction]); ok {
			engine.Analyze(d)
		}
	}
}

func processFile(filename string, selectedTypes map[string]bool, useTopologicalSort, genNoOp bool, noOpDir string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %v", filename, err)
	}

	pkg := node.Name.Name
	fmt.Printf("\n=== Analyzing file: %s ===\n", filename)

	// Create analysis engines for selected types
	engines := make(map[string]interface{})

	if selectedTypes["structs"] {
		structVisitor := NewGenericVisitor(
			NewStructNodeVisitor(fset, pkg),
			NewStructResultCollector(),
			&StructValidator{},
		)

		var structSorter ItemSorter[GoStruct]
		if useTopologicalSort {
			structSorter = NewDependencySorter(
				&StructDependencyExtractor{},
				&StructTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&StructDependencyExtractor{},
					&StructTypeNameProvider{},
				),
			)
		} else {
			structSorter = NewDependencySorter(
				&StructDependencyExtractor{},
				&StructTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&StructTypeNameProvider{},
				),
			)
		}

		structFormatter := NewGenericFormatter(
			&StructItemRenderer{},
			&SimpleOutputFormatter[GoStruct]{},
		)

		structEngine := NewAnalysisEngine(
			structVisitor,
			structSorter,
			structFormatter,
			nil,
		)

		engines["structs"] = structEngine
	}

	if selectedTypes["interfaces"] {
		interfaceVisitor := NewGenericVisitor(
			NewInterfaceNodeVisitor(fset, pkg),
			NewInterfaceResultCollector(),
			&InterfaceValidator{},
		)

		var interfaceSorter ItemSorter[GoInterface]
		if useTopologicalSort {
			interfaceSorter = NewDependencySorter(
				&InterfaceDependencyExtractor{},
				&InterfaceTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&InterfaceDependencyExtractor{},
					&InterfaceTypeNameProvider{},
				),
			)
		} else {
			interfaceSorter = NewDependencySorter(
				&InterfaceDependencyExtractor{},
				&InterfaceTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&InterfaceTypeNameProvider{},
				),
			)
		}

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

		engines["interfaces"] = interfaceEngine
	}

	if selectedTypes["functions"] {
		functionVisitor := NewGenericVisitor(
			NewFunctionNodeVisitor(fset, pkg),
			NewFunctionResultCollector(),
			&FunctionValidator{},
		)

		var functionSorter ItemSorter[GoFunction]
		if useTopologicalSort {
			functionSorter = NewDependencySorter(
				&FunctionDependencyExtractor{},
				&FunctionTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&FunctionDependencyExtractor{},
					&FunctionTypeNameProvider{},
				),
			)
		} else {
			functionSorter = NewDependencySorter(
				&FunctionDependencyExtractor{},
				&FunctionTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&FunctionTypeNameProvider{},
				),
			)
		}

		functionFormatter := NewGenericFormatter(
			&FunctionItemRenderer{},
			&SimpleOutputFormatter[GoFunction]{},
		)

		functionEngine := NewAnalysisEngine(
			functionVisitor,
			functionSorter,
			functionFormatter,
			nil,
		)

		engines["functions"] = functionEngine
	}

	if selectedTypes["variables"] {
		variableVisitor := NewGenericVisitor(
			NewVariableNodeVisitor(fset, pkg),
			NewVariableResultCollector(),
			&VariableValidator{},
		)

		var variableSorter ItemSorter[GoVariable]
		if useTopologicalSort {
			variableSorter = NewDependencySorter(
				&VariableDependencyExtractor{},
				&VariableTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&VariableDependencyExtractor{},
					&VariableTypeNameProvider{},
				),
			)
		} else {
			variableSorter = NewDependencySorter(
				&VariableDependencyExtractor{},
				&VariableTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&VariableTypeNameProvider{},
				),
			)
		}

		variableFormatter := NewGenericFormatter(
			&VariableItemRenderer{},
			&SimpleOutputFormatter[GoVariable]{},
		)

		variableEngine := NewAnalysisEngine(
			variableVisitor,
			variableSorter,
			variableFormatter,
			nil,
		)

		engines["variables"] = variableEngine
	}

	if selectedTypes["constants"] {
		constantVisitor := NewGenericVisitor(
			NewConstantNodeVisitor(fset, pkg),
			NewConstantResultCollector(),
			&ConstantValidator{},
		)

		var constantSorter ItemSorter[GoConstant]
		if useTopologicalSort {
			constantSorter = NewDependencySorter(
				&ConstantDependencyExtractor{},
				&ConstantTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&ConstantDependencyExtractor{},
					&ConstantTypeNameProvider{},
				),
			)
		} else {
			constantSorter = NewDependencySorter(
				&ConstantDependencyExtractor{},
				&ConstantTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&ConstantTypeNameProvider{},
				),
			)
		}

		constantFormatter := NewGenericFormatter(
			&ConstantItemRenderer{},
			&SimpleOutputFormatter[GoConstant]{},
		)

		constantEngine := NewAnalysisEngine(
			constantVisitor,
			constantSorter,
			constantFormatter,
			nil,
		)

		engines["constants"] = constantEngine
	}

	if selectedTypes["imports"] {
		importVisitor := NewGenericVisitor(
			NewImportNodeVisitor(fset),
			NewImportResultCollector(),
			&ImportValidator{},
		)

		var importSorter ItemSorter[GoImport]
		if useTopologicalSort {
			importSorter = NewDependencySorter(
				&ImportDependencyExtractor{},
				&ImportTypeNameProvider{},
				NewTopologicalDependencyResolver(
					&ImportDependencyExtractor{},
					&ImportTypeNameProvider{},
				),
			)
		} else {
			importSorter = NewDependencySorter(
				&ImportDependencyExtractor{},
				&ImportTypeNameProvider{},
				NewAlphabeticalDependencyResolver(
					&ImportTypeNameProvider{},
				),
			)
		}

		importFormatter := NewGenericFormatter(
			&ImportItemRenderer{},
			&SimpleOutputFormatter[GoImport]{},
		)

		importEngine := NewAnalysisEngine(
			importVisitor,
			importSorter,
			importFormatter,
			nil,
		)

		engines["imports"] = importEngine
	}

	// Analyze declarations
	for _, decl := range node.Decls {
		analyzeDecl(decl, engines)
	}

	// Print results for each selected type
	if engine, ok := engines["structs"].(*AnalysisEngine[GoStruct]); ok {
		fmt.Println("\n--- Structs (Dependency Order) ---")
		engine.PrintResults()
	}

	if engine, ok := engines["interfaces"].(*AnalysisEngine[GoInterface]); ok {
		fmt.Println("\n--- Interfaces (Dependency Order) ---")
		engine.PrintResults()

		// Generate NoOp file if requested
		if genNoOp && noOpDir != "" {
			baseFilename := filepath.Base(filename)
			noOpFilename := filepath.Join(noOpDir, "noop_"+strings.TrimSuffix(baseFilename, ".go")+"_interfaces.go")
			if err := engine.GenerateCodeFile(noOpFilename); err != nil {
				log.Printf("Failed to generate NoOp file %s: %v", noOpFilename, err)
			} else {
				fmt.Printf("Generated NoOp implementations: %s\n", noOpFilename)
			}
		}
	}

	if engine, ok := engines["functions"].(*AnalysisEngine[GoFunction]); ok {
		fmt.Println("\n--- Functions (Dependency Order) ---")
		engine.PrintResults()
	}

	if engine, ok := engines["variables"].(*AnalysisEngine[GoVariable]); ok {
		fmt.Println("\n--- Variables (Dependency Order) ---")
		engine.PrintResults()
	}

	if engine, ok := engines["constants"].(*AnalysisEngine[GoConstant]); ok {
		fmt.Println("\n--- Constants (Dependency Order) ---")
		engine.PrintResults()
	}

	if engine, ok := engines["imports"].(*AnalysisEngine[GoImport]); ok {
		fmt.Println("\n--- Imports (Dependency Order) ---")
		engine.PrintResults()
	}

	return nil
}

func walkDirectory(dir string, selectedTypes map[string]bool, useTopologicalSort, genNoOp bool, noOpDir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			return processFile(path, selectedTypes, useTopologicalSort, genNoOp, noOpDir)
		}

		return nil
	})
}

func main() {
	var (
		dirs        = flag.String("dirs", ".", "Comma-separated list of directories to analyze")
		showStructs = flag.Bool("structs", false, "Show structs")
		showIfaces  = flag.Bool("interfaces", false, "Show interfaces")
		showFuncs   = flag.Bool("functions", false, "Show functions")
		showVars    = flag.Bool("variables", false, "Show variables")
		showConsts  = flag.Bool("constants", false, "Show constants")
		showImports = flag.Bool("imports", false, "Show imports")
		showAll     = flag.Bool("all", false, "Show all types")
		topoSort    = flag.Bool("topo", true, "Use topological sorting based on dependencies")
		alphaSort   = flag.Bool("alpha", false, "Use alphabetical sorting instead of topological")
		genNoOp     = flag.Bool("noop", false, "Generate NoOp implementations for interfaces")
		noOpDir     = flag.String("noop-dir", "./noop", "Directory to save NoOp implementations")
	)

	flag.Parse()

	useTopologicalSort := *topoSort && !*alphaSort

	selectedTypes := make(map[string]bool)

	if *showAll {
		selectedTypes["structs"] = true
		selectedTypes["interfaces"] = true
		selectedTypes["functions"] = true
		selectedTypes["variables"] = true
		selectedTypes["constants"] = true
		selectedTypes["imports"] = true
	} else {
		selectedTypes["structs"] = *showStructs
		selectedTypes["interfaces"] = *showIfaces
		selectedTypes["functions"] = *showFuncs
		selectedTypes["variables"] = *showVars
		selectedTypes["constants"] = *showConsts
		selectedTypes["imports"] = *showImports
	}

	// If no specific type is selected, show all
	hasSelection := false
	for _, selected := range selectedTypes {
		if selected {
			hasSelection = true
			break
		}
	}

	if !hasSelection {
		selectedTypes["structs"] = true
		selectedTypes["interfaces"] = true
		selectedTypes["functions"] = true
		selectedTypes["variables"] = true
		selectedTypes["constants"] = true
		selectedTypes["imports"] = true
	}

	// Create NoOp output directory if needed
	if *genNoOp && *noOpDir != "" {
		if err := os.MkdirAll(*noOpDir, 0755); err != nil {
			log.Fatalf("Failed to create NoOp directory %s: %v", *noOpDir, err)
		}
	}

	directories := strings.Split(*dirs, ",")
	sort.Strings(directories)

	sortType := "Topological"
	if !useTopologicalSort {
		sortType = "Alphabetical"
	}
	fmt.Printf("Using %s sorting", sortType)
	if *genNoOp {
		fmt.Printf(" with NoOp generation enabled (output: %s)", *noOpDir)
	}
	fmt.Println()

	for _, dir := range directories {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}

		fmt.Printf("\n=== Analyzing directory: %s ===\n", dir)

		if err := walkDirectory(dir, selectedTypes, useTopologicalSort, *genNoOp, *noOpDir); err != nil {
			log.Printf("Error analyzing directory %s: %v", dir, err)
		}
	}
}
