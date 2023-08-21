import { useUserStore } from "./UserStore";

// https://stackoverflow.com/questions/41103360/
// https://stackoverflow.com/questions/23314806/
export function Agent<T>({
  jsonRequest = false,
  data,
  url,
}: {
  jsonRequest?: boolean;
  data?: any;
  url: string;
}): Promise<T> {
  let body: string;
  if (data) {
    body = jsonRequest
      ? JSON.stringify(data)
      : new URLSearchParams(data).toString();
  }

  const ReloadMessage = "invalid session, please re-login.";

  const store = useUserStore();

  return new Promise((resolve, reject) => {
    return fetch(url, {
      method: "POST",
      cache: "no-cache",
      body: body,
      headers: {
        "X-Requested-With": "XMLHttpRequest",
        "Content-Type": jsonRequest
          ? "application/json; charset=utf-8"
          : "application/x-www-form-urlencoded",
        Authorization: "Bearer " + store.user.token,
      },
    })
      .then((response) => {
        if (response.status === 400 || response.status === 403) {
          return response.json().then((e) => {
            if (e.message.includes("security requirements failed")) {
              store.Logout();
              throw new Error(ReloadMessage);
            }
          });
        }
        if (response.status === 401) {
          store.Logout();
          const message = response.headers.get("WWW-Authenticate")
            ? "Incorrect username or password."
            : ReloadMessage;
          throw new Error(message);
        }
        if (response.ok) {
          return resolve(response.json() as Promise<T>);
        }
        throw new Error("Network Error");
      })
      .catch((error) => {
        alert(error);
        if (error.message === ReloadMessage) {
          location.reload();
        } else {
          reject(error);
        }
      });
  });
}
