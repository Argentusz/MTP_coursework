import {editor, languages} from "monaco-editor";

export const mtpOperators = [
    "skip", "mov", "add", "adc", "sub", "sbb", "mul", "div", "rmd",
    "imov", "iadd", "iadc", "isub", "isbb", "imul", "idiv",
    "addf", "subf", "mulf", "divf",
    "shl", "shr", "sar", "and", "or", "xor", "not",
    "jmp", "jif", "jnf", "jiz", "jnz", "lbl", "call", "ret", "hlt", "ei", "di", "int"
]
export const mtpMacros = ["#mov32", "#break", "#start"]
export const mtpPresetters = ["@label", "@ilabel"]
export const mtpRegisters = [
    "rb0", "rb1", "rb2", "rb3", "rb4", "rb5", "rb6", "rb7",
    "rw24", "rw25", "rw26", "rw27", "rw28", "rw29", "rw30", "rw31",
    "rw16", "rw17", "rw18", "rw19", "rw20", "rw21", "rw22", "rw23",
    "rw8", "rw9", "rw10", "rw11", "rw12", "rw13", "rw14", "rw15",
    "rw0", "rw1", "rw2", "rw3", "rw4", "rw5", "rw6", "rw7",
    "rx0", "rx1", "rx2", "rx3", "rx4", "rx5", "rx6", "rx7",
    "rh0", "rh1", "rh2", "rh3", "rh4", "rh5", "rh6", "rh7",
    "rl0", "rl1", "rl2", "rl3", "rl4", "rl5", "rl6", "rl7",
]
export const mtpBuiltInAliases = [
    "$define",
    "$main",
    "$signone", "$sigfpe", "$sigtrace", "$sigsegv", "$sigterm", "$sigint", "$sigiie", "$sigill",
    "$m8", "$m16", "$m32"
]

const mtpOperatorsRegExp = new RegExp(`${mtpOperators.join("|")}`, "i")
const mtpMacrosRegExp = new RegExp(`${mtpMacros.join("|")}`, "i")
const mtpPresettersRegExp = new RegExp(`${mtpPresetters.map(el => el.replace("@", "[@]")).join("|")}`, "i")
const mtpRegistersRegExp = new RegExp(`${mtpRegisters.join("|")}`, "i")
const mtpAddressRegExp = new RegExp(`\\[(${mtpRegisters.join("|")})\\]`, "i")

const mtpTokens = [
    [mtpOperatorsRegExp, "operator"],
    [mtpMacrosRegExp, "macro"],
    [mtpPresettersRegExp, "presetter"],
    [mtpRegistersRegExp, "register"],
    [mtpAddressRegExp, "address"],
    [/\$[a-zA-Z0-9_]*( |$)/i, "alias"],
    [/0( |$)/i, "number"],
    [/[1-9][0-9]*( |$)/i, "number"],
    [/0x[0-9a-fA-F]*( |$)/i, "number"],
    [/0o[0-7]*( |$)/i, "number"],
    [/0b[01]*( |$)/i, "number"],
    [/; *TODO.*$/, "todo"],
    [/;.*$/i, "comment"],
]

const fillSuggestion = (token, kind) => ({
    label: token,
    insertText: `${token} `,
    kind,
    insertTextRules: languages.CompletionItemInsertTextRule.InsertAsSnippet,
})
const getSuggestions = () => {
    const suggestions = []

    suggestions.push(...mtpOperators.map(token => fillSuggestion(token, languages.CompletionItemKind.Operator)))
    suggestions.push(...mtpMacros.map(token => fillSuggestion(token, languages.CompletionItemKind.Keyword)))
    suggestions.push(...mtpPresetters.map(token => fillSuggestion(token, languages.CompletionItemKind.Keyword)))
    suggestions.push(...mtpRegisters.map(token => fillSuggestion(token, languages.CompletionItemKind.Variable)))
    suggestions.push(...mtpRegisters.map(token => fillSuggestion(`[${token}]`, languages.CompletionItemKind.Variable)))
    suggestions.push(...mtpBuiltInAliases.map(token => fillSuggestion(token, languages.CompletionItemKind.Unit)))

    return { provideCompletionItems: () => ({suggestions}) }
}


export const mtpLanguageDefinition = {
    tokenizer: {
        root: mtpTokens
    }
}

export const mtpColorRules = [
    { token: "operator", foreground: "#F86EA8" },
    { token: "macro", foreground: "#FFA245" },
    { token: "presetter", foreground: "#75C2B3" },
    { token: "register", foreground: "#B37EEE" },
    { token: "address", foreground: "#49B0CE" },
    { token: "alias", foreground: "#DBB8FF" },
    { token: "number", foreground: "#D7C781" },
    { token: "todo", foreground: "#A8C023" },
    { token: "comment", foreground: "#AAAAAA" },
]

export const mtpTheme = {
    base: "vs-dark",
    inherit: true,
    rules: mtpColorRules,
    colors: {}
}

export const InitMtpEditor = () => {
    languages.register({id: "mtp"})
    languages.setMonarchTokensProvider("mtp", mtpLanguageDefinition)
    languages.registerCompletionItemProvider("mtp", getSuggestions())
    editor.defineTheme("mtpTheme", mtpTheme)
}