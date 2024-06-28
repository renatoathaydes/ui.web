/**
 * Error response to executing a command in the backend.
 */
export type CommandResponseError = {
    error: string
}

/**
 * A successful response to executing a command in the backend.
 */
export type CommandResponseSuccess = {
    /** The result of the command. */
    value: any,
    /** A frontend command to execute when this value is received. */
    feCmd?: string,
}

/**
 * The result of executing a command in the backend.
 */
export type CommandResponse =
    CommandResponseSuccess | CommandResponseError

declare global {

    /**
     * Evaluate JS code in the browser.
     * 
     * If this call is triggered by a backend-provided command, then a value
     * may be sent to the feCmd that is being executed (see [CommandResponseSuccess]).
     * 
     * @param cmd command to execute
     * @param me component that triggered this command
     * @param value result of a previous call, usually a backend-provided value
     */
    export function evalWith(cmd: string, me?: any, value?: any): Promise<any>

}
