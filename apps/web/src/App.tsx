//import { useState } from 'react'
import './App.css';
import { Box, Button, HStack, Text } from '@chakra-ui/react';

// type TaxYear = {
//   year: number;
//   status: 'draft' | 'completed';
// };

function App() {
  //const [appName, setAppName] = useState('Sebas Dian app');

  //TODO: que es la verificacion estatica de typeScript
  //TypeScript analiza el codigo antes de ejecutarlo para detectar operacion incompatibles
  // const salary: number = 6_2234_32423;
  // salary.toUppercase();

  // const TaxtYear: TaxYear = {
  //   year: 'djjdsgjfsdag',
  //   status: 'pendding',
  // };

  return (
    <>
      <Box bg="tomato" w="100%" p="4" color="white">
        <h1>Tax assistan Colombia</h1>
      </Box>
      <Text>Prepare and understand your Colombia tax return</Text>
      <Text textStyle="7xl">Sprint 0</Text>
      <HStack>
        <Button>Start tax profile</Button>
      </HStack>
    </>
  );
}

export default App;
