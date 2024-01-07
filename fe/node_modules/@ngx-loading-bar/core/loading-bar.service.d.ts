import { NgZone } from '@angular/core';
import { LoadingBarState } from './loading-bar.state';
import { Observable, OperatorFunction } from 'rxjs';
import { LoadingBarConfig } from './loading-bar.config';
import * as i0 from "@angular/core";
export declare class LoadingBarService {
    private platformId;
    private config;
    private zone?;
    private readonly refs;
    private readonly streams$;
    readonly value$: Observable<number>;
    /** @deprecated use `value$` instead. */
    get progress$(): Observable<number>;
    constructor(platformId: Object, config?: LoadingBarConfig, zone?: NgZone);
    /** @deprecated use `useRef` instead. */
    start(initialValue?: number): void;
    /** @deprecated use `useRef` instead. */
    set(value: number): void;
    /** @deprecated use `useRef` instead. */
    increment(value?: number): void;
    /** @deprecated use `useRef` instead. */
    complete(): void;
    /** @deprecated use `useRef` instead. */
    stop(): void;
    useRef(id?: string): LoadingBarState;
    static ɵfac: i0.ɵɵFactoryDeclaration<LoadingBarService, [null, { optional: true; }, { optional: true; }]>;
    static ɵprov: i0.ɵɵInjectableDeclaration<LoadingBarService>;
}
export declare function runInZone<T>(zone: NgZone): OperatorFunction<T, T>;
