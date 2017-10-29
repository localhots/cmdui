export function api(path) {
    let proto = window.location.protocol,
        host = window.location.host;
    return proto + "//" + host + "/api" + path;
}

export function httpGET(url, success, error) {
    let xhr = new XMLHttpRequest();
    xhr.addEventListener("load", (e) => {
        if (xhr.status >= 400) {
            error("Request failed: " + xhr.statusText);
        } else {
            success(xhr.status, xhr.responseText);
        }
    });
    xhr.addEventListener("error", (e) => {
        error("Request failed");
    });
    xhr.addEventListener("abort", (e) => {
        error("Connection closed");
    });

    let async = true;
    xhr.open("GET", url, async);
    xhr.send(null);
    return xhr;
}

export function httpStreamGET(url, progress, complete, error) {
    let xhr = new XMLHttpRequest();
    xhr.responseType = "text";
    var lastIndex = 0;
    xhr.onreadystatechange = () => {
        let state = xhr.readyState;
        if (state === xhr.LOADING) {
            let curIndex = xhr.responseText.length;
            if (curIndex === lastIndex) {
                // No progress was made
                return;
            }

            let text = xhr.responseText.slice(lastIndex, curIndex);
            lastIndex = curIndex;
            progress(text);
        } else if (state === xhr.DONE) {
            if (xhr.status >= 400) {
                error("Request failed: " + xhr.statusText);
            } else {
                complete(xhr.status);
            }
        }
        // Ignoring states: UNSENT, OPENED, HEADERS_RECEIVED
    };
    xhr.onerror = (e) => {
        error("Request failed");
    };
    xhr.onabort = (e) => {
        error("Connection closed");
    };

    let async = true;
    xhr.open("GET", url, async);
    xhr.send(null);
    return xhr;
}

export function httpPOST(url, form, success, error) {
    let xhr = new XMLHttpRequest();
    xhr.responseType = "text";
    xhr.onreadystatechange = () => {
        let state = xhr.readyState;
        if (state === xhr.DONE) {
            if (xhr.status >= 400) {
                error("Request failed: " + xhr.statusText);
            } else {
                success(xhr.responseText);
            }
        }
        // Ignoring states: UNSENT, OPENED, HEADERS_RECEIVED, PROGRESS
    };
    xhr.onerror = (e) => {
        error("Request failed");
    };
    xhr.onabort = (e) => {
        error("Connection closed");
    };

    let async = true;
    xhr.open("POST", url, async);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.send(form);
    return xhr;
}
