import { authHeader } from '../_helpers';
import { AuthClient } from "auth_grpc_web_pb";
import { LoginRequest, RegistrationRequest } from "auth_pb";
export const userService = {
    connectGRPCAuthServer,
    gRPCLogin,
    gRPCRegistration,
    login,
    logout,
    register,
    getAll,
    getById,
    update,
    delete: _delete
};

function connectGRPCAuthServer() {
    this.client = new AuthClient("http://Host:8080", null, null);
}

async function gRPCLogin(username, password) {
    let loginRequest = new LoginRequest();
    loginRequest.setUsername(username);
    loginRequest.setPassword(password);
    let rep = "fail";
    return new Promise((resolve, reject) => {
        this.client.login(loginRequest, {}, (err, response) => {
            rep = response.toObject().token;
            console.log("in user.gRPCLogin::rep",rep);
            if (rep === "ok") {
                console.log("gRPCLogin::in promise::success")
                let myUser = {
                    userId: '',
                    username,
                    password
                };
                resolve(myUser);
            } else {
                console.log("gRPCLogin::in promise::fail")
                reject(rep);
            }
        });
    });
}

async function gRPCRegistration(username, password) {
    let registrationRequest = new RegistrationRequest();
    registrationRequest.setUsername(username);
    registrationRequest.setPassword(password);
    let rep = "fail";
    console.log("user.service::In gRPCRegistration");
    return new Promise((resolve, reject) => {
        this.client.registration(registrationRequest, {}, (err, response) => {
            rep = response.toObject().token;
            if (rep === "success") {
                console.log("gRPCRegistration::in promise::success")
                resolve(rep);
            } else {
                console.log("gRPCRegistration::in promise::fail")
                reject(rep);
            }
        });
    });
}

function login(username, password) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    };

    return fetch(`http://localhost:4000/users/authenticate`, requestOptions)
        .then(handleResponse)
        .then(user => {
            // login successful if there's a jwt token in the response
            if (user.token) {
                // store user details and jwt token in local storage to keep user logged in between page refreshes
                localStorage.setItem('user', JSON.stringify(user));
            }

            return user;
        });
}

function logout() {
    // remove user from local storage to log user out
    localStorage.removeItem('user');
}

function register(user) {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    };

    return fetch(`http://localhost:4000/users/register`, requestOptions).then(handleResponse);
}

function getAll() {
    const requestOptions = {
        method: 'GET',
        headers: authHeader()
    };

    return fetch(`http://localhost:4000/users`, requestOptions).then(handleResponse);
}


function getById(id) {
    const requestOptions = {
        method: 'GET',
        headers: authHeader()
    };

    return fetch(`http://localhost:4000/users/${id}`, requestOptions).then(handleResponse);
}

function update(user) {
    const requestOptions = {
        method: 'PUT',
        headers: { ...authHeader(), 'Content-Type': 'application/json' },
        body: JSON.stringify(user)
    };

    return fetch(`http://localhost:4000/users/${user.id}`, requestOptions).then(handleResponse);
}

// prefixed function name with underscore because delete is a reserved word in javascript
function _delete(id) {
    const requestOptions = {
        method: 'DELETE',
        headers: authHeader()
    };

    return fetch(`http://localhost:4000/users/${id}`, requestOptions).then(handleResponse);
}

function handleResponse(response) {
    return response.text().then(text => {
        const data = text && JSON.parse(text);
        if (!response.ok) {
            if (response.status === 401) {
                // auto logout if 401 response returned from api
                logout();
                location.reload(true);
            }

            const error = (data && data.message) || response.statusText;
            return Promise.reject(error);
        }

        return data;
    });
}