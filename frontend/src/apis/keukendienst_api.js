import axios from 'axios';

const HOST = `http://${window.location.host}`

export const getKeukendienst = () => axios.get(`${HOST}/api/keukendienst`)
export const setKeukendienst = (info) => axios.put(`${HOST}/api/keukendienst`, info)
