const questions = [
  "У вас есть аллергия?",
  "Вы принимаете постоянные медикаменты?",
  "Есть ли у вас хронические заболевания?",
  "Были ли у вас операции за последний год?",
  "Есть ли у вас диабет?",
  "Имеются ли проблемы с сердцем?",
  "Вы курите?",
  "Вы употребляете алкоголь регулярно?",
  "У вас повышенное давление?",
  "Бывают ли у вас обмороки?",
  "Вы беременны (если применимо)?",
  "Вы испытываете стресс регулярно?",
  "Есть ли у вас заболевания печени?",
  "Есть ли у вас заболевания почек?",
  "Есть ли у вас инфекции в настоящее время?"
];

const form = document.getElementById('questionForm');

questions.forEach((question, index) => {
  const div = document.createElement('div');
  div.style.marginTop = '15px';

  const label = document.createElement('label');
  label.textContent = question;

  const checkbox = document.createElement('input');
  checkbox.type = 'checkbox';
  checkbox.name = `question_${index}`;
  checkbox.value = 'yes';

  div.appendChild(label);
  div.appendChild(checkbox);
  form.appendChild(div);
});

document.getElementById('submitBtn').addEventListener('click', function (e) {
  e.preventDefault();

  const answers = {};
  questions.forEach((_, index) => {
    const checkbox = document.querySelector(`input[name="question_${index}"]`);
    answers[`question_${index}`] = checkbox.checked ? 'Да' : 'Нет';
  });

  localStorage.setItem('answers', JSON.stringify(answers));

  // Переход на страницу "Спасибо"
  window.location.href = 'success.html';
});
