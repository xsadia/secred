import type { NextPage } from "next";
import { Text } from "@chakra-ui/react";
import { useContext, useEffect } from "react";
import { AuthContext } from "../context/AuthContext";
import Router from "next/router";

const Home: NextPage = () => {
  const { user, isAuthenticated, isLoading } = useContext(AuthContext);
  // useEffect(() => {
  //   if (!isLoading && !isAuthenticated) {
  //     Router.push("/login");
  //   }
  // }, [isAuthenticated, isLoading]);
  return (
    <div>
      <Text fontSize="6xl">Hello {user?.username}</Text>
    </div>
  );
};

export default Home;
