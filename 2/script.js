function sendDataToServer() {
  const formData = {
    firstName: localStorage.getItem("firstName"),
    lastName: localStorage.getItem("lastName"),
    middleName: localStorage.getItem("middleName"),
    birthDate: localStorage.getItem("birthDate"),
    phone: localStorage.getItem("phone"),
    gender: localStorage.getItem("gender"),
    answers: JSON.parse(localStorage.getItem("answers"))
  };

  fetch('http://localhost:8080/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  })
  .then(res => res.blob())
  .then(blob => {
    const link = document.createElement('a');
    link.href = window.URL.createObjectURL(blob);
    link.download = 'documents.zip';
    link.click();
  })
  .catch(err => console.error(err));
}
