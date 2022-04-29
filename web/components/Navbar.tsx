import { Button, Flex, Heading } from "@chakra-ui/react";

type User = {
  email: string;
  username: string;
};

type NavbarProps = {
  isAuthenticated?: boolean;
  user?: User;
};

export const Navbar = ({ isAuthenticated, user }: NavbarProps) => {
  return (
    <Flex
      as="header"
      position="fixed"
      py="12px"
      px="6px"
      w="100vw"
      align="center"
      borderBottom="1px"
      borderColor="gray.300"
    >
      <Flex justifyContent="space-between" w="100%">
        <Heading>Secred</Heading>
        {isAuthenticated && <Button>Logout</Button>}
      </Flex>
    </Flex>
  );
};
