import React,{useState,useEffect} from 'react';

function DownloadForm() {
  const [file, setFile]=useState([])
  useEffect(() => {
    fetchFile();
  }, [])

  const [password, setPassword]=useState([])

  useEffect(() => {
    console.log(file)
  }, [file])

  const fetchFile=async()=>{
    const response=await fetch('http://localhost:10000/api/v1/file/' + password);
    setFile(response)
    console.log(password)
    console.log(file)    
  }

  const handleSubmit=async(event)=>{
    this.preventDefault();
    const data = new FormData(event.target);
    console.log(this.state.password)
    console.log(data)
  }
 const handlePasswordChange=async(event) => {
   console.log("Updating password", password)
    setPassword(event.target.value);
 }

  return (
    <form onSubmit={handleSubmit} encType = "multipart/form-data" >
           
        <label htmlFor="file">Enter passphrase</label>
        <input id="text" name="password" type="text" value={password}
        onChange={handlePasswordChange} />

        <button>Download file!</button>
      </form>
  );
}

export default DownloadForm;