package lsp

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/calico32/goose/ast"
	std_platform "github.com/calico32/goose/lib/std/platform"
	"github.com/calico32/goose/parser"
	"github.com/calico32/goose/scanner"
	"github.com/calico32/goose/token"
	"github.com/calico32/goose/validator"
	"go.lsp.dev/jsonrpc2"
	. "go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func (ls *LanguageServer) Location(fset *token.FileSet, node ast.Node) Location {
	start := fset.Position(node.Pos())
	end := fset.Position(node.End())
	return Location{
		URI: ls.uris[start.Filename],
		Range: Range{
			Start: Position{
				Line:      uint32(start.Line) - 1,
				Character: uint32(start.Column) - 1,
			},
			End: Position{
				Line:      uint32(end.Line) - 1,
				Character: uint32(end.Column) - 1,
			},
		},
	}
}

func (ls *LanguageServer) Range(fset *token.FileSet, node ast.Node) Range {
	start := fset.Position(node.Pos())
	end := fset.Position(node.End())
	return Range{
		Start: Position{
			Line:      uint32(start.Line) - 1,
			Character: uint32(start.Column) - 1,
		},
		End: Position{
			Line:      uint32(end.Line) - 1,
			Character: uint32(end.Column) - 1,
		},
	}
}
func (ls *LanguageServer) ensureInitialized() *jsonrpc2.Error {
	if ls.state == ServerStateIdle {
		err := jsonrpc2.Errorf(jsonrpc2.ServerNotInitialized, "server not initialized")
		ls.logger.Sugar().Error(err)
		return err
	}
	return nil
}

func (ls *LanguageServer) Initialize(ctx context.Context, params *InitializeParams) (result *InitializeResult, err error) {
	if ls.state != ServerStateIdle {
		ls.logger.Sugar().Errorf("Initialize called in invalid state: %s", ls.state)
	}
	ls.state = ServerStateRunning
	return &InitializeResult{
		ServerInfo: &ServerInfo{
			Name:    "goose-lsp",
			Version: std_platform.Version,
		},
		Capabilities: ServerCapabilities{
			TextDocumentSync:       TextDocumentSyncKindIncremental,
			HoverProvider:          true,
			DocumentSymbolProvider: true,
			DefinitionProvider:     true,
			DeclarationProvider:    true,
		},
	}, nil
}

func (ls *LanguageServer) Initialized(ctx context.Context, params *InitializedParams) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}
	return
}

func (ls *LanguageServer) SetTrace(ctx context.Context, params *SetTraceParams) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}
	return
}

func (ls *LanguageServer) LogTrace(ctx context.Context, params *LogTraceParams) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}
	return
}

func (ls *LanguageServer) Shutdown(ctx context.Context) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}
	return
}

func (ls *LanguageServer) Exit(ctx context.Context) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}
	err = ls.conn.Close()
	if err != nil {
		return
	}
	<-ls.conn.Done()
	os.Exit(0)
	return
}

func (ls *LanguageServer) CodeAction(ctx context.Context, params *CodeActionParams) (result []CodeAction, err error) {
	err = notImplemented("CodeAction")
	return
}

func (ls *LanguageServer) CodeLens(ctx context.Context, params *CodeLensParams) (result []CodeLens, err error) {
	err = notImplemented("CodeLens")
	return
}

func (ls *LanguageServer) CodeLensRefresh(ctx context.Context) (err error) {
	err = notImplemented("CodeLensRefresh")
	return
}

func (ls *LanguageServer) CodeLensResolve(ctx context.Context, params *CodeLens) (result *CodeLens, err error) {
	err = notImplemented("CodeLensResolve")
	return
}

func (ls *LanguageServer) ColorPresentation(ctx context.Context, params *ColorPresentationParams) (result []ColorPresentation, err error) {
	err = notImplemented("ColorPresentation")
	return
}

func (ls *LanguageServer) Completion(ctx context.Context, params *CompletionParams) (result *CompletionList, err error) {
	err = notImplemented("Completion")
	return
}

func (ls *LanguageServer) CompletionResolve(ctx context.Context, params *CompletionItem) (result *CompletionItem, err error) {
	err = notImplemented("CompletionResolve")
	return
}

func (ls *LanguageServer) Declaration(ctx context.Context, params *DeclarationParams) (result []Location /* Declaration | DeclarationLink[] | null */, err error) {
	ls.logger.Sugar().Debugf("Declaration: %s", params.TextDocument.URI.Filename())
	moduleMu := ls.modules[params.TextDocument.URI]
	if moduleMu == nil {
		ls.logger.Sugar().Errorf("module not found: %s", params.TextDocument.URI.Filename())
		return nil, errors.New("module not found")
	}

	fset, ok := ls.fsets[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("fileset not found: %s", params.TextDocument.URI.Filename())
	}

	position := fset.GetPosition("file:"+params.TextDocument.URI.Filename(), int(params.Position.Line)+1, int(params.Position.Character)+1)
	pos := fset.Pos(position)
	if pos == token.NoPos {
		ls.logger.Sugar().Errorf("position not found: %#v", position)
		return []Location{}, nil
	}

	module := moduleMu.Lock()
	defer moduleMu.Unlock()

	node := module.FindNode(pos)
	if node == nil {
		ls.logger.Sugar().Errorf("node not found: %s", pos)
		return []Location{}, nil
	}

	// TODO: find the declaration of the node
	return []Location{ls.Location(fset, node)}, nil
}

func (ls *LanguageServer) Definition(ctx context.Context, params *DefinitionParams) (result []Location /* Definition | DefinitionLink[] | null */, err error) {
	return ls.Declaration(ctx, &DeclarationParams{
		TextDocumentPositionParams: params.TextDocumentPositionParams,
	})
}

// CmpPositions compares two positions and returns:
//
//	-1 if a < b
//	0 if a == b
//	1 if a > b
func CmpPositions(a, b Position) int {
	if a.Line < b.Line {
		return -1
	}
	if a.Line > b.Line {
		return 1
	}
	if a.Character < b.Character {
		return -1
	}
	if a.Character > b.Character {
		return 1
	}
	return 0
}

func (ls *LanguageServer) DidChange(ctx context.Context, params *DidChangeTextDocumentParams) (err error) {
	ls.logger.Sugar().Debugf("DidChange: %s", params.TextDocument.URI.Filename())
	ls.logger.Sugar().Debugf("changes: %#v", params.ContentChanges)
	sourceMu, ok := ls.sourceFiles[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("document not found: %s", params.TextDocument.URI.Filename())
		return errors.New("document not found")
	}
	source := sourceMu.Lock()
	before := string(source)

	fset, ok := ls.fsets[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("fileset not found: %s", params.TextDocument.URI.Filename())
	}
	file := fset.File("file:" + params.TextDocument.URI.Filename())

	lineOffsets := file.Lines()

	for _, change := range params.ContentChanges {
		start := change.Range.Start
		end := change.Range.End
		startOffset := lineOffsets[int(start.Line)] + int(start.Character)
		endOffset := len(source) // end of the file (because end.Line could be 1 greater than the number of lines)
		if end.Line < uint32(len(lineOffsets)) {
			endOffset = lineOffsets[int(end.Line)] + int(end.Character)
		}

		ls.logger.Sugar().Debugf("start: %d, end: %d, source: %d", startOffset, endOffset, len(source))
		source = append(source[:startOffset], append([]byte(change.Text), source[endOffset:]...)...)
		ls.logger.Sugar().Debugf("source: %s", source)

		// if startOffset > endOffset {
		// 	if len(change.Text) > 0 {
		// 		ls.logger.Sugar().Errorf("range start > end and has content: %#v", change)
		// 	} else {
		// 		// start > end, delete the range
		// 		source = append(source[:endOffset], source[startOffset-1:]...)
		// 	}
		// } else {
		// if start < end, this replaces the range
		// if start == end, this inserts at the position
		// }
		// } else {
		// 	source = []byte(change.Text)
		// }
	}

	if before == string(source) {
		ls.logger.Sugar().Debugf("no changes")
		return
	}

	sourceMu.Update(source)
	err = ls.checkModule(ctx, params.TextDocument.URI)
	return
}

func (ls *LanguageServer) DidChangeConfiguration(ctx context.Context, params *DidChangeConfigurationParams) (err error) {
	err = notImplemented("DidChangeConfiguration")
	return
}

func (ls *LanguageServer) DidChangeWatchedFiles(ctx context.Context, params *DidChangeWatchedFilesParams) (err error) {
	err = notImplemented("DidChangeWatchedFiles")
	return
}

func (ls *LanguageServer) DidChangeWorkspaceFolders(ctx context.Context, params *DidChangeWorkspaceFoldersParams) (err error) {
	err = notImplemented("DidChangeWorkspaceFolders")
	return
}

func (ls *LanguageServer) DidClose(ctx context.Context, params *DidCloseTextDocumentParams) (err error) {
	err = notImplemented("DidClose")
	return
}

func (ls *LanguageServer) DidCreateFiles(ctx context.Context, params *CreateFilesParams) (err error) {
	err = notImplemented("DidCreateFiles")
	return
}

func (ls *LanguageServer) DidDeleteFiles(ctx context.Context, params *DeleteFilesParams) (err error) {
	err = notImplemented("DidDeleteFiles")
	return
}

func (ls *LanguageServer) DidOpen(ctx context.Context, params *DidOpenTextDocumentParams) (err error) {
	if err := ls.ensureInitialized(); err != nil {
		return err
	}

	if params.TextDocument.LanguageID != "goose" && !strings.HasSuffix(params.TextDocument.URI.Filename(), ".goose") {
		ls.logger.Sugar().Errorf("ignoring document with languageID %s", params.TextDocument.LanguageID)
		return
	}

	ls.sourceFiles[params.TextDocument.URI] = &Mutexed[[]byte]{
		v: []byte(params.TextDocument.Text),
	}
	err = ls.checkModule(ctx, params.TextDocument.URI)
	return
}

func (ls *LanguageServer) checkModule(ctx context.Context, documentUri DocumentURI) error {
	ls.logger.Sugar().Debugf("checking module: %s", documentUri.Filename())
	sourceMu, ok := ls.sourceFiles[documentUri]
	if !ok {
		ls.logger.Sugar().Errorf("document not found: %s", documentUri.Filename())
		return errors.New("document not found")
	}
	ls.uris[documentUri.Filename()] = documentUri
	src := sourceMu.Lock()
	defer sourceMu.Unlock()

	fset := token.NewFileSet()
	ls.fsets[documentUri] = fset
	t := ls.timer.Mark("parse")
	module, err := parser.ParseFile(fset, "file:"+documentUri.Filename(), src, nil)
	t.Done()
	if err != nil {
		if errs, ok := err.(scanner.ErrorList); ok {
			t := ls.timer.Mark("writeDiagnostics")
			ls.logger.Sugar().Debugf("parsed %s with %d errors", documentUri.Filename(), len(errs))
			// report parse errors to the client
			diagnostics := map[DocumentURI][]Diagnostic{}
			for _, e := range errs {
				doc, ok := ls.uris[strings.TrimPrefix(e.Pos.Filename, "file:")]
				if !ok {
					doc = uri.File(strings.TrimPrefix(e.Pos.Filename, "file:"))
				}
				if diagnostics[doc] == nil {
					diagnostics[doc] = []Diagnostic{}
				}
				diagnostics[doc] = append(diagnostics[doc], Diagnostic{
					Message:  e.Msg,
					Severity: DiagnosticSeverityError,
					Range: Range{
						Start: Position{
							Line:      uint32(e.Pos.Line) - 1,
							Character: uint32(e.Pos.Column) - 1,
						},
						End: Position{
							Line:      uint32(e.Pos.Line) - 1,
							Character: uint32(e.Pos.Column) - 1,
						},
					},
				})
			}
			t.Done()

			t = ls.timer.Mark("publishDiagnostics")
			for uri, diags := range diagnostics {
				ls.parseErrors[uri] = diags
				err = ls.client.PublishDiagnostics(ctx, &PublishDiagnosticsParams{
					URI:         uri,
					Diagnostics: diags,
				})
				if err != nil {
					return err
				}
			}
			t.Done()
		}
		return err
	}

	// parse OK!
	ls.parseErrors[documentUri] = nil
	diagnostics := []Diagnostic{}

	t = ls.timer.Mark("validate")
	// move onto validation
	v, err := validator.New(module, fset, false, nil, nil, nil)
	if err != nil {
		return err
	}
	_, err = v.Check()
	if err != nil {
		return err
	}
	t.Done()

	t = ls.timer.Mark("writeValidationDiagnostics")
	problems := v.Diagnostics()
	if len(problems) > 0 {
		for _, problem := range problems {
			diagnostics = append(diagnostics, Diagnostic{
				Message:  problem.Message,
				Severity: problem.Severity,
				Range:    ls.Range(fset, problem.Node),
				Source:   "goose",
			})
		}
	}
	t.Done()

	t = ls.timer.Mark("publishDiagnostics")
	ls.client.PublishDiagnostics(ctx, &PublishDiagnosticsParams{
		URI:         documentUri,
		Diagnostics: diagnostics,
	})
	t.Done()

	ls.modules[documentUri] = &Mutexed[*ast.Module]{
		v: module,
	}
	return nil
}

func (ls *LanguageServer) DidRenameFiles(ctx context.Context, params *RenameFilesParams) (err error) {
	err = notImplemented("DidRenameFiles")
	return
}

func (ls *LanguageServer) DidSave(ctx context.Context, params *DidSaveTextDocumentParams) (err error) {
	err = notImplemented("DidSave")
	return
}

func (ls *LanguageServer) DocumentColor(ctx context.Context, params *DocumentColorParams) (result []ColorInformation, err error) {
	err = notImplemented("DocumentColor")
	return
}

func (ls *LanguageServer) DocumentHighlight(ctx context.Context, params *DocumentHighlightParams) (result []DocumentHighlight, err error) {
	err = notImplemented("DocumentHighlight")
	return
}

func (ls *LanguageServer) DocumentLink(ctx context.Context, params *DocumentLinkParams) (result []DocumentLink, err error) {
	err = notImplemented("DocumentLink")
	return
}

func (ls *LanguageServer) DocumentLinkResolve(ctx context.Context, params *DocumentLink) (result *DocumentLink, err error) {
	err = notImplemented("DocumentLinkResolve")
	return
}

func (ls *LanguageServer) DocumentSymbol(ctx context.Context, params *DocumentSymbolParams) (result []any /* []SymbolInformation | []DocumentSymbol */, err error) {
	return []any{}, nil

	moduleMu, ok := ls.modules[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Debugf("module not found: %s", params.TextDocument.URI.Filename())
		// probably a byproduct of a parse error, don't report it to the client
		return nil, nil
	}

	if parseErrors, ok := ls.parseErrors[params.TextDocument.URI]; ok {
		if parseErrors != nil {
			// don't try to get symbols if there are parse errors
			// (AST is likely incomplete or incorrect)
			return nil, nil
		}
	}

	module := moduleMu.Lock()
	defer moduleMu.Unlock()

	fset, ok := ls.fsets[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("fileset not found: %s", params.TextDocument.URI.Filename())
	}

	symbols := []any{}
	for _, stmt := range module.Stmts {
		switch stmt := stmt.(type) {
		case *ast.ConstStmt:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Ident.Name,
				Kind:           SymbolKindConstant,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Ident),
			})
		case *ast.ExprStmt:
			switch expr := stmt.X.(type) {
			case *ast.FuncExpr:
				symbols = append(symbols, DocumentSymbol{
					Name:           expr.Name.Name,
					Kind:           SymbolKindFunction,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, expr.Name),
				})
			}
		case *ast.StructStmt:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Name.Name,
				Kind:           SymbolKindStruct,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Name),
			})
		case *ast.NativeConst:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Ident.Name,
				Kind:           SymbolKindConstant,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Ident),
			})
		case *ast.NativeFunc:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Name.Name,
				Kind:           SymbolKindFunction,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Name),
			})
		case *ast.NativeStruct:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Name.Name,
				Kind:           SymbolKindStruct,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Name),
			})
		case *ast.LetStmt:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Ident.Name,
				Kind:           SymbolKindVariable,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Ident),
			})
		case *ast.SymbolStmt:
			symbols = append(symbols, DocumentSymbol{
				Name:           stmt.Ident.Name,
				Kind:           SymbolKindKey,
				Range:          ls.Range(fset, stmt),
				SelectionRange: ls.Range(fset, stmt.Ident),
			})
		case *ast.ExportDeclStmt:
			switch decl := stmt.Stmt.(type) {
			case *ast.ConstStmt:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Ident.Name,
					Kind:           SymbolKindConstant,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Ident),
				})
			case *ast.ExprStmt:
				if fn, ok := decl.X.(*ast.FuncExpr); ok {
					symbols = append(symbols, DocumentSymbol{
						Name:           fn.Name.Name,
						Kind:           SymbolKindFunction,
						Range:          ls.Range(fset, stmt),
						SelectionRange: ls.Range(fset, fn.Name),
					})
				}
			case *ast.StructStmt:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Name.Name,
					Kind:           SymbolKindStruct,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Name),
				})
			case *ast.LetStmt:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Ident.Name,
					Kind:           SymbolKindVariable,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Ident),
				})
			case *ast.NativeConst:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Ident.Name,
					Kind:           SymbolKindConstant,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Ident),
				})
			case *ast.NativeFunc:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Name.Name,
					Kind:           SymbolKindFunction,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Name),
				})
			case *ast.NativeStruct:
				symbols = append(symbols, DocumentSymbol{
					Name:           decl.Name.Name,
					Kind:           SymbolKindStruct,
					Range:          ls.Range(fset, stmt),
					SelectionRange: ls.Range(fset, decl.Name),
				})
			}

		}
	}

	ls.logger.Sugar().Debugf("symbols: %#v", symbols)

	return symbols, nil
}

func (ls *LanguageServer) ExecuteCommand(ctx context.Context, params *ExecuteCommandParams) (result interface{}, err error) {
	err = notImplemented("ExecuteCommand")
	return
}

func (ls *LanguageServer) FoldingRanges(ctx context.Context, params *FoldingRangeParams) (result []FoldingRange, err error) {
	err = notImplemented("FoldingRanges")
	return
}

func (ls *LanguageServer) Formatting(ctx context.Context, params *DocumentFormattingParams) (result []TextEdit, err error) {
	err = notImplemented("Formatting")
	return
}

func (ls *LanguageServer) Hover(ctx context.Context, params *HoverParams) (result *Hover, err error) {
	if err := ls.ensureInitialized(); err != nil {
		return nil, err
	}

	moduleMu, ok := ls.modules[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("module not found: %s", params.TextDocument.URI.Filename())
		return nil, errors.New("module not found")
	}

	module := moduleMu.Lock()
	defer moduleMu.Unlock()

	fset, ok := ls.fsets[params.TextDocument.URI]
	if !ok {
		ls.logger.Sugar().Errorf("fileset not found: %s", params.TextDocument.URI.Filename())
	}

	position := fset.GetPosition("file:"+params.TextDocument.URI.Filename(), int(params.Position.Line)+1, int(params.Position.Character)+1)
	pos := fset.Pos(position)
	if pos == token.NoPos {
		ls.logger.Sugar().Errorf("position not found: %#v", position)
		return &Hover{}, nil
	}

	node := module.FindNode(pos)
	if node == nil {
		ls.logger.Sugar().Errorf("node not found: %s", pos)
		return &Hover{}, nil
	}

	// TODO: find the declaration of the node

	return &Hover{}, nil
}

func (ls *LanguageServer) Implementation(ctx context.Context, params *ImplementationParams) (result []Location, err error) {
	err = notImplemented("Implementation")
	return
}

func (ls *LanguageServer) IncomingCalls(ctx context.Context, params *CallHierarchyIncomingCallsParams) (result []CallHierarchyIncomingCall, err error) {
	err = notImplemented("IncomingCalls")
	return
}

func (ls *LanguageServer) LinkedEditingRange(ctx context.Context, params *LinkedEditingRangeParams) (result *LinkedEditingRanges, err error) {
	err = notImplemented("LinkedEditingRange")
	return
}

func (ls *LanguageServer) Moniker(ctx context.Context, params *MonikerParams) (result []Moniker, err error) {
	err = notImplemented("Moniker")
	return
}

func (ls *LanguageServer) OnTypeFormatting(ctx context.Context, params *DocumentOnTypeFormattingParams) (result []TextEdit, err error) {
	err = notImplemented("OnTypeFormatting")
	return
}

func (ls *LanguageServer) OutgoingCalls(ctx context.Context, params *CallHierarchyOutgoingCallsParams) (result []CallHierarchyOutgoingCall, err error) {
	err = notImplemented("OutgoingCalls")
	return
}

func (ls *LanguageServer) PrepareCallHierarchy(ctx context.Context, params *CallHierarchyPrepareParams) (result []CallHierarchyItem, err error) {
	err = notImplemented("PrepareCallHierarchy")
	return
}

func (ls *LanguageServer) PrepareRename(ctx context.Context, params *PrepareRenameParams) (result *Range, err error) {
	err = notImplemented("PrepareRename")
	return
}

func (ls *LanguageServer) RangeFormatting(ctx context.Context, params *DocumentRangeFormattingParams) (result []TextEdit, err error) {
	err = notImplemented("RangeFormatting")
	return
}

func (ls *LanguageServer) References(ctx context.Context, params *ReferenceParams) (result []Location, err error) {
	err = notImplemented("References")
	return
}

func (ls *LanguageServer) Rename(ctx context.Context, params *RenameParams) (result *WorkspaceEdit, err error) {
	err = notImplemented("Rename")
	return
}

func (ls *LanguageServer) Request(ctx context.Context, method string, params interface{}) (result interface{}, err error) {
	err = notImplemented("Request")
	return
}

func (ls *LanguageServer) SemanticTokensFull(ctx context.Context, params *SemanticTokensParams) (result *SemanticTokens, err error) {
	err = notImplemented("SemanticTokensFull")
	return
}

func (ls *LanguageServer) SemanticTokensFullDelta(ctx context.Context, params *SemanticTokensDeltaParams) (result interface{} /* SemanticTokens | SemanticTokensDelta */, err error) {
	err = notImplemented("SemanticTokensFullDelta")
	return
}

func (ls *LanguageServer) SemanticTokensRange(ctx context.Context, params *SemanticTokensRangeParams) (result *SemanticTokens, err error) {
	err = notImplemented("SemanticTokensRange")
	return
}

func (ls *LanguageServer) SemanticTokensRefresh(ctx context.Context) (err error) {
	err = notImplemented("SemanticTokensRefresh")
	return
}

func (ls *LanguageServer) ShowDocument(ctx context.Context, params *ShowDocumentParams) (result *ShowDocumentResult, err error) {
	err = notImplemented("ShowDocument")
	return
}

func (ls *LanguageServer) SignatureHelp(ctx context.Context, params *SignatureHelpParams) (result *SignatureHelp, err error) {
	err = notImplemented("SignatureHelp")
	return
}

func (ls *LanguageServer) Symbols(ctx context.Context, params *WorkspaceSymbolParams) (result []SymbolInformation, err error) {
	err = notImplemented("Symbols")
	return
}

func (ls *LanguageServer) TypeDefinition(ctx context.Context, params *TypeDefinitionParams) (result []Location, err error) {
	err = notImplemented("TypeDefinition")
	return
}

func (ls *LanguageServer) WillCreateFiles(ctx context.Context, params *CreateFilesParams) (result *WorkspaceEdit, err error) {
	err = notImplemented("WillCreateFiles")
	return
}

func (ls *LanguageServer) WillDeleteFiles(ctx context.Context, params *DeleteFilesParams) (result *WorkspaceEdit, err error) {
	err = notImplemented("WillDeleteFiles")
	return
}

func (ls *LanguageServer) WillRenameFiles(ctx context.Context, params *RenameFilesParams) (result *WorkspaceEdit, err error) {
	err = notImplemented("WillRenameFiles")
	return
}

func (ls *LanguageServer) WillSave(ctx context.Context, params *WillSaveTextDocumentParams) (err error) {
	err = notImplemented("WillSave")
	return
}

func (ls *LanguageServer) WillSaveWaitUntil(ctx context.Context, params *WillSaveTextDocumentParams) (result []TextEdit, err error) {
	err = notImplemented("WillSaveWaitUntil")
	return
}

func (ls *LanguageServer) WorkDoneProgressCancel(ctx context.Context, params *WorkDoneProgressCancelParams) (err error) {
	err = notImplemented("WorkDoneProgressCancel")
	return
}

func notImplemented(method string) error {
	return jsonrpc2.Errorf(jsonrpc2.MethodNotFound, "method %q not implemented", method)
}
