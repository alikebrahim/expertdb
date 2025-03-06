import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export interface Expert {
  id: string;
  name: string;
  affiliation: string;
  is_bahraini: boolean;
  isced_field_id: string;
  is_available: boolean;
  bio: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    name: string;
    email: string;
    role: string;
  };
}

export interface IscedField {
  id: string;
  name: string;
}

export const login = async (email: string, password: string): Promise<AuthResponse> => {
  const response = await api.post("/api/auth/login", { email, password });
  return response.data;
};

export const getExperts = async (filters: {
  name?: string;
  affiliation?: string;
  is_bahraini?: boolean;
  isced_field_id?: string;
  is_available?: boolean;
  page?: number;
  limit?: number;
} = {}): Promise<Expert[]> => {
  const params = new URLSearchParams();
  if (filters.name) params.append("name", filters.name);
  if (filters.affiliation) params.append("affiliation", filters.affiliation);
  if (filters.is_bahraini !== undefined) params.append("is_bahraini", String(filters.is_bahraini));
  if (filters.isced_field_id) params.append("isced_field_id", filters.isced_field_id);
  if (filters.is_available !== undefined) params.append("is_available", String(filters.is_available));
  if (filters.page) params.append("page", String(filters.page));
  if (filters.limit) params.append("limit", String(filters.limit));

  const response = await api.get("/api/experts", { params });
  return response.data;
};

export const getIscedFields = async (): Promise<IscedField[]> => {
  const response = await api.get("/api/isced/fields");
  return response.data;
};

export default api;
