import axios from "axios";
import { ElMessage } from "element-plus";
import { useStore } from "@/stores";

export const useAxios = axios.create({
  baseURL: "/api",
});

export interface baseResponse<T> {
  code: number;
  data: T;
  msg: string;
}
export interface listDataType<T> {
  count: number;
  list: T[];
}

useAxios.interceptors.request.use((config) => {
  const store = useStore();
  const token = store.userInfo.token;
  console.log(token);
  // config.headers["Content-Type"] = "multipart/form-data";
  config.headers["Authorization"] = "Bearer " + token;
  return config;
});
useAxios.interceptors.response.use(
  (response) => {
    if (response.status !== 200) {
      // 失败的
      console.log("服务失败", response.status);
      ElMessage.error(response.statusText);
      return Promise.reject(response.statusText);
    }
    const store = useStore();
    // 获取 token
    const newToken = response.headers["x-jwt-token"];
    store.setToken(newToken);
    return response.data;
  },
  (err) => {
    console.log("服务错误", err);
    ElMessage.error(err.message);
    return Promise.reject(err.message);
  },
);
