import { FormEventHandler, ReactNode } from "react";
import { Flex } from "@chakra-ui/react";

type FormLayoutProps = {
  onSubmit: FormEventHandler;
  children: ReactNode;
};

export const FormLayout = ({ children, onSubmit }: FormLayoutProps) => {
  return (
    <Flex
      align="center"
      justifyContent="center"
      height="100vh"
      background="gray.200"
    >
      <Flex
        as="form"
        direction="column"
        background="gray.300"
        p={12}
        rounded={6}
        onSubmit={onSubmit}
      >
        {children}
      </Flex>
    </Flex>
  );
};
