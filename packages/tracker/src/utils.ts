export function monkeypatchPost<T extends Record<K, (...args: any[]) => any>, K extends keyof T>(
    self: T,
    funcName: K,
    callback: (...args: Parameters<T[K]>) => void) {
    const orig = self[funcName];
    self[funcName] = ((...args: Parameters<typeof orig>): ReturnType<typeof orig> => {
        const result = orig.apply(self, args)
        callback(...args)
        return result;
    }) as T[K]
}
