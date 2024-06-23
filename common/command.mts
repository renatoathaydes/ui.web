export type CommandResponseError = {
    error: string
}

export type CommandResponseSuccess = {
    value: any,
    feCmd?: string,
}

export type CommandResponse =
    CommandResponseSuccess | CommandResponseError
