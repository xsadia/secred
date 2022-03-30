import { createContext, ReactNode, useEffect, useState } from "react";
import { setCookie, parseCookies, destroyCookie } from "nookies";
import { useRouter } from "next/router";

type AuthContextData = {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  signIn: (credentials: Credentials) => Promise<void>;
};

type Credentials = {
  email: string;
  password: string;
};

type User = {
  username: string;
  email: string;
};

const TOKEN_EXPIRED_ERROR = "Token is expired";

export const AuthContext = createContext({} as AuthContextData);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const isAuthenticated = !!user;
  const [isLoading, setIsloading] = useState<boolean>(true);
  const router = useRouter();

  useEffect(() => {
    const { "@secred:token": token } = parseCookies();

    if (token) {
      (async () => {
        const response = await fetch("http://localhost:1337/user/me", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `bearer ${token}`,
          },
        });

        const data = await response.json();

        if (response.status !== 200) {
          if (data.error === TOKEN_EXPIRED_ERROR) {
            destroyCookie(undefined, "@secred:token");
          }

          router.push("/login");
        }

        setUser({
          email: data.email,
          username: data.username,
        });
      })();

      setIsloading(false);
      return;
    }

    router.push("/login");
  }, []);

  const signIn = async ({ email, password }: Credentials) => {
    const response = await fetch("http://localhost:1337/auth", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (response.status !== 200) {
      throw new Error(data.error);
    }

    setCookie(undefined, "@secred:token", data.token, {
      maxAge: 60 * 60 * 24 * 7,
    });

    setUser({
      email: data.user.email,
      username: data.user.username,
    });
  };

  const logout = () => {
    destroyCookie(undefined, "@secred:token");
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isLoading, user, signIn }}>
      {children}
    </AuthContext.Provider>
  );
};
