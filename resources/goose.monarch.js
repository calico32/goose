/** @type {import('monaco-editor').languages.IMonarchLanguage} */
const gooseMonarch = {
  defaultToken: 'invalid',

  keywords: [
    'let',
    'const',
    'memo',
    'fn',
    'end',
    'return',
    'if',
    'else',
    'repeat',
    'while',
    'times',
    'forever',
    'break',
    'continue',
    'for',
    'in',
  ],

  typeKeywords: [],

  // prettier-ignore
  operators: [
    '=',
    '+', '-', '*', '/', '%', '**',
    '+=', '-=', '*=', '/=', '%=', '**=',
    '++', '--',
    '>', '<', '>=', '<=',
    '==', '!=',
    '&&', '||', '!'
  ],

  langconstants: ['true', 'false', 'null', '_'],

  symbols: /[=><!~?:&|+\-*/^%]+/,
  escapes: /\\(?:[abfnrtv\\"'$]|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,
  ident: /[a-zA-Z_$][a-zA-Z0-9_$]*/,
  upperIdent: /[A-Z][a-zA-Z0-9_$]*/,
  digits: /\d+(_\d+)*/,
  octaldigits: /[0-7]+(_[0-7]+)*/,
  binarydigits: /[01]+(_[01]+)*/,
  hexdigits: /[[0-9a-fA-F]+(_[0-9a-fA-F]+)*/,

  tokenizer: {
    root: [
      // function decl
      [
        /(fn)(\s+)([a-zA-Z_$][a-zA-Z0-9_$]*)(\()/,
        [
          { token: 'keyword' },
          { token: 'brace' },
          { token: 'function' },
          { token: 'paren', next: '@params' },
        ],
      ],

      // numbers
      [/(@digits)[eE]([-+]?(@digits))?/, 'number'],
      [/(@digits)\.(@digits)([eE][-+]?(@digits))?/, 'number'],
      [/0x(@hexdigits)/, 'number'],
      [/0o?(@octaldigits)/, 'number'],
      [/0b(@binarydigits)/, 'number'],
      [/(@digits)/, 'number'],

      // function call
      [/(@ident)(?=\()/, 'function'],

      // identifiers and keywords
      [/@upperIdent/, 'class'], // to show class names nicely
      [
        /@ident/,
        {
          cases: {
            '@langconstants': 'lang-constant',
            '@keywords': 'keyword',
            '@default': 'identifier',
          },
        },
      ],

      // whitespace
      { include: '@whitespace' },

      // delimiters and operators
      [/[{}()[\]]/, 'brace'],
      [
        /@symbols/,
        {
          cases: {
            '@operators': 'operator',
            '@default': '',
          },
        },
      ],

      // delimiter: after number because of .\d floats
      [/[;,.]/, 'delimiter'],

      // strings
      [/"([^"\\]|\\.)*$/, 'string-invalid'], // non-teminated string
      [/"/, { token: 'string', bracket: '@open', next: '@string' }],
    ],

    params: [
      [/@ident/, 'parameter'],
      [/,/, 'delimiter'],
      [/\)/, { token: 'paren', next: '@pop' }],
    ],

    whitespace: [
      [/[ \t\r\n]+/, 'white'],
      [/\/\*/, 'comment', '@comment'],
      [/\/\/.*$/, 'comment'],
    ],

    comment: [
      [/[^/*]+/, 'comment'],
      ['\\*/', 'comment', '@pop'],
      [/[/*]/, 'comment'],
    ],

    string: [
      [/\$\{/, { token: 'string-interpolated', next: '@string_interpolated' }],
      [/\$[a-zA-Z_][\w$]*/, 'string-interpolated'],
      [/[^\\"$]+/, 'string'],
      [/@escapes/, 'string-escape'],
      [/\\./, 'string-escape-invalid'],
      [/"/, { token: 'string', bracket: '@close', next: '@pop' }],
    ],

    string_interpolated: [
      [/\{/, { token: 'string-interpolated', next: '@string_interpolated' }],
      [/\}/, { token: 'string-interpolated', next: '@pop' }],
      { include: 'root' },
    ],
  },
}
