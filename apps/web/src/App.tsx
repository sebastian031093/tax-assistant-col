import { useState } from 'react';
import './App.css';

// type TaxYear = {
//   year: number;
//   status: 'draft' | 'completed';
// };

function App() {
  const [appName, setAppName] = useState('Sebas Dian app');

  //TODO: que es la verificacion estatica de typeScript
  //TypeScript analiza el codigo antes de ejecutarlo para detectar operacion incompatibles
  // const salary: number = 6_2234_32423;
  // salary.toUppercase();

  // const TaxtYear: TaxYear = {
  //   year: 'djjdsgjfsdag',
  //   status: 'pendding',
  // };

  return (
    <div>
      <h1>Hi from taxes Colombia App. {appName}</h1>
    </div>
  );
}

export default App;
