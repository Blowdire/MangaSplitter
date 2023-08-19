import { useEffect, useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet, ChooseFile } from '../wailsjs/go/main/App';

function App() {
  const [resultText, setResultText] = useState('Please enter your name below 👇');
  const [name, setName] = useState('');
  const updateName = (e) => setName(e.target.value);
  const [currentProgress, setcurrentProgress] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const updateResultText = (result) => setResultText(result);
  

  useEffect(() => {
    window.runtime.EventsOn('pageDone', (currentProgress) => {
      console.log(currentProgress)
      const {PageNumber, CurrentTotalPages} = currentProgress;
      if(CurrentTotalPages !== totalPages){
        setTotalPages(CurrentTotalPages)
      }
      setcurrentProgress(PageNumber)
    })
  }, []);
  

 function greet() {
    ChooseFile();
  }
  return (
    <div id="App">
      <div id="result" className="result">{resultText}</div>
      <div id="input" className="input-box">
        <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" />
        <button className="btn" onClick={greet}>Greet</button>
      </div>
    <span>{currentProgress}/{totalPages}</span>
    </div>
  );
}

export default App;
