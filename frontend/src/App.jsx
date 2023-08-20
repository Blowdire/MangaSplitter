import { useEffect, useState } from 'react';
import './App.css';
import { ChooseFile } from '../wailsjs/go/main/App';
import { Button, Col, Progress, Row, Slider } from 'antd';
import './style.css';
function App() {
  const [currentProgress, setcurrentProgress] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [sliderValue, setSliderValue] = useState(25);

  useEffect(() => {
    window.runtime.EventsOn('pageDone', (currentProgress) => {
      console.log(currentProgress)
      const { PageNumber, CurrentTotalPages } = currentProgress;
      if (CurrentTotalPages !== totalPages) {
        setTotalPages(CurrentTotalPages)
      }
      setcurrentProgress(Math.floor(PageNumber / CurrentTotalPages * 100))
    })
  }, []);
  const handleSliderChange = (value) => {
    setSliderValue(value);
  };

  function ProcessFile() {
    ChooseFile(sliderValue);
  }
  return (
    <div id="App">
      <Row justify={'center'} align='middle' style={{ marginTop: 10 }}>
        <Col><span style={{ fontSize: 32 }}>Manga splitter</span></Col>
      </Row>
      <Row justify={'center'} align='middle' style={{ marginTop: 10 }}>
        <Col><span style={{ fontSize: 18, }}>Choose a manga pdf file in order to split double pages</span></Col>
      </Row>
      <Row justify={'center'} align='middle' style={{ marginTop: 10 }}>
        <Col><Button onClick={() => ProcessFile()}>Choose pdf file</Button></Col>
      </Row>
      <Row style={{ margin: '20px', textAlign: 'center' }}>
        <Col span={24}>
          <Row justify={'center'}> <span style={{ fontSize: 18 }}>Select the quality level(higher quality leads to bigger file)</span></Row>
          <Row style={{ width: '100%' }}>
            <Slider
              style={{ width: '100%' }}
              min={1}
              max={100}
              value={sliderValue}
              onChange={handleSliderChange}
            />
          </Row>
        </Col>

        <Row style={{ marginTop: '20px', width:'100%' }} justify={'center'}>
         <Col><span style={{ fontSize: 18 }}>Quality: {sliderValue}</span></Col> 
        </Row>
      </Row>
      <Row justify={'center'} align='middle' style={{ marginTop: 10, marginLeft: 15, marginRight: 15, color: 'white' }}>
        <Progress status={currentProgress === 100 ? 'success' : 'active'} percent={currentProgress} strokeColor={{ '0%': '#108ee9', '100%': '#87d068' }} style={{ color: 'white' }} />
      </Row></div>
  );
}

export default App;
