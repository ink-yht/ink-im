export interface tokenPayload {
    exp: number;
    Uid: number;
}
export const parseToken = (token: string): tokenPayload => {
    let payload = token.split(".")[1];
    return JSON.parse(
        decodeURIComponent(
            escape(window.atob(payload.replace(/-/g, "+").replace(/_/g, "/"))),
        ),
    );
};
