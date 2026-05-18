import axios from "axios";

const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

export default api;
