import { CommandResponse } from '../../common/command.mjs';

export class BackendCommandRunner {
    verbose: boolean = false;

    async runScript(js: string): Promise<CommandResponse> {
        console.log('Running js', js);
        try {
            const value = await evalWith(js, this);
            console.log(`Success: ${value}`);
            return this.asCommandResponse(value);
        } catch (e) {
            console.warn(e);
            return { error: e.toString() };
        }
    }

    private asCommandResponse(value: any): CommandResponse {
        if (typeof value === 'object') {
            if ("error" in value) {
                return value as CommandResponse;
            }
            if ("value" in value) {
                return value as CommandResponse;
            }
        }
        return { value };
    }
}
