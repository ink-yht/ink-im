import { type baseResponse, useAxios } from "@/api/index";

export interface LoginRequest {
  email: string;
  password: string;
}

export function LoginApi(data: LoginRequest): Promise<baseResponse<string>> {
  return useAxios.post("/users/login", data);
}

// export function UserInfoApi() {
//   return useAxios.post("/users/info", data);
// }
