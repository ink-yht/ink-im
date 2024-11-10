import { type baseResponse, useAxios } from "@/api/index";

// LoginApi 用户登录
export interface LoginRequest {
  email: string;
  password: string;
}

export function LoginApi(data: LoginRequest): Promise<baseResponse<string>> {
  return useAxios.post("/users/login", data);
}

// LoginApi 用户信息
export interface VerificationQuestionType {
  problem1?: string;
  problem2?: string;
  problem3?: string;
  answer1?: string;
  answer2?: string;
  answer3?: string;
}

export interface UserInfoType {
  id: number;
  email: string;
  phone: string;
  nickname: string;
  abstract: string;
  avatar: string;
  recallMessage?: string;
  friendOnline: boolean;
  sound: boolean;
  secureLink: boolean;
  savePwd: boolean;
  searchUser: number;
  verification: number;
  verificationQuestion: VerificationQuestionType;
}

export function UserInfoApi(): Promise<baseResponse<UserInfoType>> {
  return useAxios.get("/users/info");
}
