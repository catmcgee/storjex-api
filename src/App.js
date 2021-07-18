import logo from './logo.svg';
import './App.css';
import UploadForm from './components/uploadForm'
import DownloadForm from './components/downloadForm'

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <UploadForm></UploadForm>
        <DownloadForm></DownloadForm>
      </header>
    </div>
  );
}

export default App;
