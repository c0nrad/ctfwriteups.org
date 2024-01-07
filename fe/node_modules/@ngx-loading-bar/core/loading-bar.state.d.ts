import { Observable } from 'rxjs';
import { LoadingBarConfig } from './loading-bar.config';
export declare class LoadingBarState {
    private config;
    private state;
    private requests;
    private disabled;
    private readonly stream$;
    private _value$;
    constructor(config?: LoadingBarConfig);
    get value$(): Observable<number>;
    start(initialValue?: number): void;
    stop(): void;
    complete(): void;
    disable(): void;
    set(value: number): void;
    increment(value?: number): void;
    private next;
    private timer$;
    private _increment;
}
