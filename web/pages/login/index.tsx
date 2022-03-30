import { FormEventHandler, useContext } from "react";
import {
  Button,
  Flex,
  FormControl,
  FormErrorMessage,
  Heading,
  Input,
  Link,
  useToast,
} from "@chakra-ui/react";
import { useFormik } from "formik";
import { FormLayout } from "../../components/FormLayout";
import { AuthContext } from "../../context/AuthContext";
import NextLink from "next/link";
import * as Yup from "yup";
import { useRouter } from "next/router";

type FormData = {
  email: string;
  password: string;
};

const LoginPage = () => {
  const { signIn } = useContext(AuthContext);
  const toast = useToast();
  const router = useRouter();

  const validationSchema = Yup.object({
    email: Yup.string()
      .required("E-mail é um campo obrigatório")
      .email("Informe um e-mail válido"),
    password: Yup.string().required("Senha é um campo obrigatório"),
  });

  const handleSignIn = async (formData: FormData) => {
    try {
      await signIn(formData);
      router.push("/");
    } catch (e) {
      toast({
        status: "error",
        title: `${e}`,
        position: "top-right",
        variant: "left-accent",
        isClosable: true,
      });
    }
  };

  const formik = useFormik({
    initialValues: {
      email: "",
      password: "",
    },
    onSubmit: async (values, actions) => {
      handleSignIn(values);

      actions.resetForm();
    },
    validationSchema: validationSchema,
  });

  return (
    <FormLayout onSubmit={formik.handleSubmit as FormEventHandler}>
      <Heading mb={6} color="gray.600">
        Log in
      </Heading>
      <FormControl
        isInvalid={!!formik.errors.email && formik.touched.email}
        mb={3}
      >
        <Input
          name="email"
          placeholder="E-mail"
          variant="outline"
          background="white"
          _placeholder={{ color: "black" }}
          onChange={formik.handleChange}
          value={formik.values.email}
          type="email"
          onBlur={formik.handleBlur}
        />
        <FormErrorMessage>{formik.errors.email}</FormErrorMessage>
      </FormControl>
      <FormControl
        mb={6}
        isInvalid={!!formik.errors.password && formik.touched.password}
      >
        <Input
          name="password"
          placeholder="Password"
          variant="outline"
          background="white"
          _placeholder={{ color: "black" }}
          onChange={formik.handleChange}
          value={formik.values.password}
          type="password"
          onBlur={formik.handleBlur}
        />
        <FormErrorMessage>{formik.errors.password}</FormErrorMessage>
      </FormControl>
      <Button colorScheme="blue" type="submit">
        Sign in
      </Button>
      <Flex mt="3" justifyContent="center">
        <NextLink href="/signup">
          <Link _hover={{ color: "blue.600" }} mr="2">
            Cadastrar conta
          </Link>
        </NextLink>
      </Flex>
    </FormLayout>
  );
};

export default LoginPage;
