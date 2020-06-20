import axios from "axios";
import {HOST} from './main'

export function APILogin(dispatch: Function, redirect: Function, code: string|null) {
    axios.get(HOST+"auth/callback/"+code).then(val => {
        dispatch('connect', val.data)
        redirect()
    });
}

export function APICheckValidity(token: string, action: Function) {
    axios.get(HOST+"auth/verify/"+token).then(val => {
            action(200, val.data)
        })
        .catch(err => {
            if (err.response && err.response.status) {
                action(err.response.status, null)
            } else {
                action(-1, null)
            }
        });
}