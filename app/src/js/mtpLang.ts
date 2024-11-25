import {editor, languages} from "monaco-editor";

export const mtpOperators = /skip|mov|add|adc|sub|sbb|mul|div|rmd|imov|iadd|iadc|isub|isbb|imul|idiv|addf|subf|mulf|divf|shl|shr|sar|and|or|xor|not|jmp|jif|jnf|jiz|jnz|lbl|call|ret|hlt|ei|di|int"/

export const mtpTokens = [
    [mtpOperators, "operator"],
    [/#[a-zA-Z0-9]*( |$)/i, "macro"],
    [/@[a-zA-Z0-9]*( |$)/i, "presetter"],
    [/\$[a-zA-Z0-9_]*( |$)/i, "alias"],
    [/\[r[xhlwb][0-9]*]/i, "address"],
    [/r[xhlwb][0-9]*( |$)/i, "register"],
    [/[1-9][0-9]*( |$)/i, "number"],
    [/0x[0-9a-fA-F]*( |$)/i, "number"],
    [/0o[0-7]*( |$)/i, "number"],
    [/0b[01]*( |$)/i, "number"],
    [/; *TODO.*$/, "todo"],
    [/;.*$/i, "comment"],
]

export const mtpLanguageDefinition = {
    tokenizer: {
        root: mtpTokens
    }
}

export const mtpColorRules = [
    { token: "operator", foreground: "#F86EA8" },
    { token: "macro", foreground: "#FFA245" },
    { token: "presetter", foreground: "#75C2B3" },
    { token: "alias", foreground: "#DBB8FF" },
    { token: "address", foreground: "#DACA77" },
    { token: "register", foreground: "#F0F0F0" },
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


languages.register({ id: "mtp" })
languages.setMonarchTokensProvider("mtp", mtpLanguageDefinition)
editor.defineTheme("mtpTheme", mtpTheme)