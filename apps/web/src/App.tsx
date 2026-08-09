import { useState } from 'react';
import './App.css';

function App() {
  const [appName, setAppName] = useState('Sebas Dian app');

  return (
    <div>
      <h1>Hi from taxes Colombia App. {appName}</h1>
    </div>
  );
}

export default App;
