import { createContext, useContext } from "react";
import type { User } from "../generated/api";

type AuthContextValue = {
  me: User | undefined;
};

export const AuthContext = createContext<AuthContextValue>({ me: undefined });

export const useAuth = () => useContext(AuthContext);
