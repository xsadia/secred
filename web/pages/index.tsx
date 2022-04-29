import type {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from "next";
import {
  Box,
  Flex,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  VStack,
} from "@chakra-ui/react";

import { parseCookies } from "nookies";
import { Navbar } from "../components/Navbar";
import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";

type Item = {
  id: string;
  name: string;
  min: number;
  max: number;
  quantity: number;
};

export default function Home({
  items,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const { user, isAuthenticated } = useContext(AuthContext);
  return (
    <>
      <Navbar isAuthenticated={isAuthenticated} />
      <Flex
        h="100vh"
        align="center"
        justifyContent="center"
        background="gray.200"
      >
        <Box w="80%" overflowY="auto" maxHeight="300px">
          <Table colorScheme="blackAlpha">
            <Thead>
              <Tr>
                <Th>Item</Th>
                <Th>Quantidade</Th>
                <Th>Min</Th>
                <Th>Max</Th>
              </Tr>
            </Thead>
            <Tbody>
              {items.map((item: Item) => (
                <Tr key={item.id}>
                  <Td>{item.name}</Td>
                  <Td>{item.quantity}</Td>
                  <Td>{item.min}</Td>
                  <Td>{item.max}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Box>
      </Flex>
    </>
  );
}

export const getServerSideProps: GetServerSideProps = async (
  ctx: GetServerSidePropsContext
) => {
  const { ["@secred:token"]: token } = parseCookies(ctx);
  if (!token) {
    return {
      redirect: {
        destination: "/login",
        permanent: false,
      },
    };
  }

  const response = await fetch("http://localhost:1337/warehouse?count=10", {
    method: "GET",
    headers: {
      Authorization: `bearer ${token}`,
    },
  });

  const data = await response.json();

  return {
    props: {
      items: data,
    },
  };
};
