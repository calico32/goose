package lsp

import (
	"context"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)



func (ls *LanguageServer) CodeAction(ctx context.Context, params *protocol.CodeActionParams) (result []protocol.CodeAction, err error) {
	err = notImplemented("CodeAction")
	return
}

func (ls *LanguageServer) CodeLens(ctx context.Context, params *protocol.CodeLensParams) (result []protocol.CodeLens, err error) {
	err = notImplemented("CodeLens")
	return
}

func (ls *LanguageServer) CodeLensRefresh(ctx context.Context) (err error) {
	err = notImplemented("CodeLensRefresh")
	return
}

func (ls *LanguageServer) CodeLensResolve(ctx context.Context, params *protocol.CodeLens) (result *protocol.CodeLens, err error) {
	err = notImplemented("CodeLensResolve")
	return
}

func (ls *LanguageServer) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) (result []protocol.ColorPresentation, err error) {
	err = notImplemented("ColorPresentation")
	return
}

func (ls *LanguageServer) Completion(ctx context.Context, params *protocol.CompletionParams) (result *protocol.CompletionList, err error) {
	err = notImplemented("Completion")
	return
}

func (ls *LanguageServer) CompletionResolve(ctx context.Context, params *protocol.CompletionItem) (result *protocol.CompletionItem, err error) {
	err = notImplemented("CompletionResolve")
	return
}

func (ls *LanguageServer) Declaration(ctx context.Context, params *protocol.DeclarationParams) (result []protocol.Location /* Declaration | DeclarationLink[] | null */, err error) {
	err = notImplemented("Declaration")
	return
}

func (ls *LanguageServer) Definition(ctx context.Context, params *protocol.DefinitionParams) (result []protocol.Location /* Definition | DefinitionLink[] | null */, err error) {
	err = notImplemented("Definition")
	return
}

func (ls *LanguageServer) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) (err error) {
	err = notImplemented("DidChange")
	return
}

func (ls *LanguageServer) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) (err error) {
	err = notImplemented("DidChangeConfiguration")
	return
}

func (ls *LanguageServer) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) (err error) {
	err = notImplemented("DidChangeWatchedFiles")
	return
}

func (ls *LanguageServer) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) (err error) {
	err = notImplemented("DidChangeWorkspaceFolders")
	return
}

func (ls *LanguageServer) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) (err error) {
	err = notImplemented("DidClose")
	return
}

func (ls *LanguageServer) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (err error) {
	err = notImplemented("DidCreateFiles")
	return
}

func (ls *LanguageServer) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (err error) {
	err = notImplemented("DidDeleteFiles")
	return
}

func (ls *LanguageServer) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) (err error) {
	err = notImplemented("DidOpen")
	return
}

func (ls *LanguageServer) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (err error) {
	err = notImplemented("DidRenameFiles")
	return
}

func (ls *LanguageServer) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) (err error) {
	err = notImplemented("DidSave")
	return
}

func (ls *LanguageServer) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) (result []protocol.ColorInformation, err error) {
	err = notImplemented("DocumentColor")
	return
}

func (ls *LanguageServer) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) (result []protocol.DocumentHighlight, err error) {
	err = notImplemented("DocumentHighlight")
	return
}

func (ls *LanguageServer) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) (result []protocol.DocumentLink, err error) {
	err = notImplemented("DocumentLink")
	return
}

func (ls *LanguageServer) DocumentLinkResolve(ctx context.Context, params *protocol.DocumentLink) (result *protocol.DocumentLink, err error) {
	err = notImplemented("DocumentLinkResolve")
	return
}

func (ls *LanguageServer) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) (result []interface{} /* []SymbolInformation | []DocumentSymbol */, err error) {
	err = notImplemented("DocumentSymbol")
	return
}

func (ls *LanguageServer) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (result interface{}, err error) {
	err = notImplemented("ExecuteCommand")
	return
}

func (ls *LanguageServer) Exit(ctx context.Context) (err error) {
	err = notImplemented("Exit")
	return
}

func (ls *LanguageServer) FoldingRanges(ctx context.Context, params *protocol.FoldingRangeParams) (result []protocol.FoldingRange, err error) {
	err = notImplemented("FoldingRanges")
	return
}

func (ls *LanguageServer) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) (result []protocol.TextEdit, err error) {
	err = notImplemented("Formatting")
	return
}

func (ls *LanguageServer) Hover(ctx context.Context, params *protocol.HoverParams) (result *protocol.Hover, err error) {
	err = notImplemented("Hover")
	return
}

func (ls *LanguageServer) Implementation(ctx context.Context, params *protocol.ImplementationParams) (result []protocol.Location, err error) {
	err = notImplemented("Implementation")
	return
}

func (ls *LanguageServer) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) (result []protocol.CallHierarchyIncomingCall, err error) {
	err = notImplemented("IncomingCalls")
	return
}

func (ls *LanguageServer) Initialized(ctx context.Context, params *protocol.InitializedParams) (err error) {
	err = notImplemented("Initialized")
	return
}

func (ls *LanguageServer) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (result *protocol.LinkedEditingRanges, err error) {
	err = notImplemented("LinkedEditingRange")
	return
}

func (ls *LanguageServer) LogTrace(ctx context.Context, params *protocol.LogTraceParams) (err error) {
	err = notImplemented("LogTrace")
	return
}

func (ls *LanguageServer) Moniker(ctx context.Context, params *protocol.MonikerParams) (result []protocol.Moniker, err error) {
	err = notImplemented("Moniker")
	return
}

func (ls *LanguageServer) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) (result []protocol.TextEdit, err error) {
	err = notImplemented("OnTypeFormatting")
	return
}

func (ls *LanguageServer) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) (result []protocol.CallHierarchyOutgoingCall, err error) {
	err = notImplemented("OutgoingCalls")
	return
}

func (ls *LanguageServer) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) (result []protocol.CallHierarchyItem, err error) {
	err = notImplemented("PrepareCallHierarchy")
	return
}

func (ls *LanguageServer) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (result *protocol.Range, err error) {
	err = notImplemented("PrepareRename")
	return
}

func (ls *LanguageServer) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) (result []protocol.TextEdit, err error) {
	err = notImplemented("RangeFormatting")
	return
}

func (ls *LanguageServer) References(ctx context.Context, params *protocol.ReferenceParams) (result []protocol.Location, err error) {
	err = notImplemented("References")
	return
}

func (ls *LanguageServer) Rename(ctx context.Context, params *protocol.RenameParams) (result *protocol.WorkspaceEdit, err error) {
	err = notImplemented("Rename")
	return
}

func (ls *LanguageServer) Request(ctx context.Context, method string, params interface{}) (result interface{}, err error) {
	err = notImplemented("Request")
	return
}

func (ls *LanguageServer) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (result *protocol.SemanticTokens, err error) {
	err = notImplemented("SemanticTokensFull")
	return
}

func (ls *LanguageServer) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (result interface{} /* SemanticTokens | SemanticTokensDelta */, err error) {
	err = notImplemented("SemanticTokensFullDelta")
	return
}

func (ls *LanguageServer) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (result *protocol.SemanticTokens, err error) {
	err = notImplemented("SemanticTokensRange")
	return
}

func (ls *LanguageServer) SemanticTokensRefresh(ctx context.Context) (err error) {
	err = notImplemented("SemanticTokensRefresh")
	return
}

func (ls *LanguageServer) SetTrace(ctx context.Context, params *protocol.SetTraceParams) (err error) {
	err = notImplemented("SetTrace")
	return
}

func (ls *LanguageServer) ShowDocument(ctx context.Context, params *protocol.ShowDocumentParams) (result *protocol.ShowDocumentResult, err error) {
	err = notImplemented("ShowDocument")
	return
}

func (ls *LanguageServer) Shutdown(ctx context.Context) (err error) {
	err = notImplemented("Shutdown")
	return
}

func (ls *LanguageServer) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (result *protocol.SignatureHelp, err error) {
	err = notImplemented("SignatureHelp")
	return
}

func (ls *LanguageServer) Symbols(ctx context.Context, params *protocol.WorkspaceSymbolParams) (result []protocol.SymbolInformation, err error) {
	err = notImplemented("Symbols")
	return
}

func (ls *LanguageServer) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) (result []protocol.Location, err error) {
	err = notImplemented("TypeDefinition")
	return
}

func (ls *LanguageServer) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (result *protocol.WorkspaceEdit, err error) {
	err = notImplemented("WillCreateFiles")
	return
}

func (ls *LanguageServer) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (result *protocol.WorkspaceEdit, err error) {
	err = notImplemented("WillDeleteFiles")
	return
}

func (ls *LanguageServer) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (result *protocol.WorkspaceEdit, err error) {
	err = notImplemented("WillRenameFiles")
	return
}

func (ls *LanguageServer) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) (err error) {
	err = notImplemented("WillSave")
	return
}

func (ls *LanguageServer) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) (result []protocol.TextEdit, err error) {
	err = notImplemented("WillSaveWaitUntil")
	return
}

func (ls *LanguageServer) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) (err error) {
	err = notImplemented("WorkDoneProgressCancel")
	return
}

func (ls *LanguageServer) Initialize(ctx context.Context, params *protocol.InitializeParams) (result *protocol.InitializeResult, err error) {
	err = notImplemented("Initialize")
	return
}

func notImplemented(method string) error {
	return jsonrpc2.Errorf(jsonrpc2.MethodNotFound, "method %q not implemented", method)
}
