/// <reference types="vite/client" />

interface ImportMetaEnv {
    /** Base URL of the API to proxy (optional - used in vite.config.ts dev proxy) */
    readonly VITE_API_URL?: string;
    /** Full URL of the endpoint to fetch on startup (optional) */
    readonly VITE_API_ENDPOINT?: string;
    /** Display name shown in the header */
    readonly VITE_APP_NAME?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
