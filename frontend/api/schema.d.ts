/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
    "/": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description index page */
        get: operations["Index"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/{filename}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description serve file */
        get: operations["ServeRootFile"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/login": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** @description login */
        post: operations["Login"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/logout": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** @description logout */
        post: operations["Logout"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/profile": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description user profile */
        get: operations["Profile"];
        /** @description update profile */
        put: operations["UpdateProfile"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/category_with_count": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description categories with unread entry count. */
        get: operations["CategoryAndUnreadEntryCount"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/category/{category_id}/entry": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description unread entries */
        get: operations["UnreadEntry"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/category": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description category list */
        get: operations["Categories"];
        put?: never;
        /** @description register category */
        post: operations["RegisterCategory"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/subscription": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description subscription list */
        get: operations["Subscriptions"];
        put?: never;
        /** @description register subscription */
        post: operations["RegisterSubscription"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/opml": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description export subscription via opml document */
        get: operations["OpmlExport"];
        put?: never;
        /** @description import opml into subscription */
        post: operations["OpmlImport"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/examine_subscription": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** @description probe web site */
        post: operations["ExamineSubscription"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/subscription/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        /** @description change subscription */
        put: operations["ChangeSubscription"];
        post?: never;
        /** @description delete subscription */
        delete: operations["DeleteSubscription"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/category/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /** @description delete category */
        delete: operations["DeleteCategory"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/pin/asread": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** @description set readflag */
        post: operations["SetAsRead"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/pin": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** @description return Pinned items */
        get: operations["PinnedItems"];
        put?: never;
        /** @description set pin into entry */
        post: operations["SetPin"];
        /** @description remove all pins */
        delete: operations["RemoveAllPin"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/update_password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** @description update password */
        post: operations["UpdatePassword"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        /** @description ログイン情報 */
        Authorization: {
            autoseen: boolean;
            token: string;
        };
        /** @description 個人設定 */
        Profile: {
            autoseen: boolean;
            /** Format: uint64 */
            entryCount: number;
            onLoginSkipPinList: boolean;
            /** Format: uint64 */
            substringLength: number;
        };
        /** @description ピン止めしたアイテム */
        PinnedItem: {
            title: string;
            url: string;
            /** Format: uint64 */
            serial: number;
            /** Format: uint64 */
            feed_id: number;
            update_at: string;
        };
        /** @description カテゴリごとに未読記事数 */
        CategoryAndUnreadEntryCount: {
            name: string;
            /** Format: uint64 */
            count: number;
            /**
             * Format: uint64
             * @description category ID
             */
            id: number;
        };
        /** @description カテゴリに属した未読記事一覧 */
        UnreadEntry: {
            /** Format: uint64 */
            serial: number;
            /** Format: uint64 */
            feed_id: number;
            title: string;
            description: string;
            /** Format: uint64 */
            date_epoch: number;
            /** @enum {string} */
            readflag: "Unseen" | "Seen" | "Setpin";
            url: string;
            /** Format: uint64 */
            subscription_id: number;
            site_title: string;
        };
        /** @description カテゴリ一覧 */
        Category: {
            /** Format: uint64 */
            id: number;
            name: string;
        };
        /** @description フィード探索 */
        ExamineSubscription: {
            success: boolean;
            title: string;
            url: string;
            preview_feed: components["schemas"]["ExamineFeed"][];
        };
        /** @description フィード探索におけるフィード詳細 */
        ExamineFeed: {
            title: string;
            url: string;
            date: string;
        };
        SimpleResult: {
            result: string;
        };
        /** @description 既読情報 */
        AsRead: {
            /** Format: uint64 */
            feed_id: number;
            /** Format: uint64 */
            serial: number;
        };
        /** @description カテゴリおよび購読一覧 */
        Subscription: {
            /**
             * Format: uint64
             * @description カテゴリID
             */
            id: number;
            /**
             * Format: string
             * @description カテゴリ名
             */
            name: string;
            /** @description カテゴリに属するフィード一覧 */
            subscription: components["schemas"]["CategorySubscription"][];
        };
        CategorySubscription: {
            /**
             * Format: uint64
             * @description フィードID
             */
            id: number;
            title: string;
            /** Format: uint64 */
            category_id: number;
            /** @description 最終アクセス時のレスポンスコード */
            http_status: string;
            /** @description フィード配信元サイトURL */
            siteurl: string;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    Index: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "text/html": string;
                };
            };
        };
    };
    ServeRootFile: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description filename */
                filename: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "text/html": string;
                };
            };
            /** @description missing file */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    Login: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    id: string;
                    password: string;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Authorization"];
                };
            };
            /** @description Error */
            default: {
                headers: {
                    /** @description error */
                    "WWW-Authenticate"?: string;
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    Logout: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": Record<string, never>;
                };
            };
        };
    };
    Profile: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Profile"];
                };
            };
        };
    };
    UpdateProfile: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": components["schemas"]["Profile"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
            /** @description error */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    CategoryAndUnreadEntryCount: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CategoryAndUnreadEntryCount"][];
                };
            };
            /** @description error */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    UnreadEntry: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description category id */
                category_id: number;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["UnreadEntry"][];
                };
            };
        };
    };
    Categories: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Category"][];
                };
            };
        };
    };
    RegisterCategory: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    name: string;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
            /** @description error */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description duplicate error */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    Subscriptions: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Subscription"][];
                };
            };
        };
    };
    RegisterSubscription: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    /** @description Site Title */
                    title: string;
                    /** @description RSS Feed URL */
                    rss: string;
                    /** @description Site URL */
                    url: string;
                    /** Format: uint64 */
                    category: number;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
            /** @description error */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description duplicate error */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    OpmlExport: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        /** @description XML document */
                        xml: string;
                    };
                };
            };
        };
    };
    OpmlImport: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    /** @description Opml xml document */
                    xml: string;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        done: boolean;
                    };
                };
            };
        };
    };
    ExamineSubscription: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    /** @description Site URL */
                    url: string;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExamineSubscription"];
                };
            };
        };
    };
    ChangeSubscription: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description subscription id */
                id: number;
            };
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    /** Format: uint64 */
                    category: number;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
        };
    };
    DeleteSubscription: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description subscription id */
                id: number;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
        };
    };
    DeleteCategory: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description category id */
                id: number;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
        };
    };
    SetAsRead: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": components["schemas"]["AsRead"][];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
            /** @description error */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    PinnedItems: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PinnedItem"][];
                };
            };
        };
    };
    SetPin: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    /** @enum {string} */
                    readflag: "Unseen" | "Seen" | "Setpin";
                    /** Format: uint64 */
                    serial: number;
                    /** Format: uint64 */
                    feed_id: number;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        /** @enum {string} */
                        readflag: "Unseen" | "Seen" | "Setpin";
                    };
                };
            };
        };
    };
    RemoveAllPin: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
        };
    };
    UpdatePassword: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": {
                    password_old: string;
                    password: string;
                    passwordc: string;
                };
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SimpleResult"];
                };
            };
        };
    };
}
