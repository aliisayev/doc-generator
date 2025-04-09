const medicalGroups = [
  {
    title: "1. Ürək-damar xəstəlikləri",
    triggerName: "heart_disease",
    followUps: [
      "Revmakardit", "Miokard infarktı", "Oynaqların kəskin revmatizmi", "Böyrək xəstəlikləri", "Stenokardiya", "Ürək catışmazlığı", "Digər"
    ],
    allowCustom: true
  },
  {
    title: "2. Arterial təzyiq",
    triggerName: "pressure",
    followUps: ["Hipertenziya", "Hiportenziya"]
  },
  {
    title: "3. Qan sisteminin xəstəlikləri",
    triggerName: "blood_diseases",
    followUps: ["Hemorragik diatez", "Anemiya", "Laxtalanma sisteminin anomaliyası", "Leykemiya", "Digər"],
    allowCustom: true
  },
  {
    title: "4. QBB xəstəlikləri",
    triggerName: "ent_diseases",
    followUps: ["Frontal Sinusit", "Maksillar Sinusit", "Tonzil Hipertroniyası", "Digər"],
    allowCustom: true
  },
  {
    title: "5. Ağciyər xəstəlikləri",
    triggerName: "lungs",
    followUps: ["Astma", "Keçirilmiş tuberkulyoz", "Aktiv tuberkulyoz", "Bronxit", "Digər"],
    allowCustom: true
  },
  {
    title: "6. Psixo-nevroloji xəstəliklər",
    triggerName: "neuro",
    followUps: ["Epilepsiya", "Nevrolgiyalar", "Paraliç", "Psixiatrik müalicə", "Digər"],
    allowCustom: true
  },
  {
    title: "7. Qaraciyər xəstəlikləri",
    triggerName: "liver",
    followUps: ["Hepatit A", "Hepatit B", "Hepatit C", "Qaraciyər catışmazlığı", "Digər"],
    allowCustom: true
  },
  {
    title: "8. Endokrin xəstəlikləri",
    triggerName: "endocrine",
    followUps: ["Diabet", "Hipotireoz", "Paratireoz", "Hipertireoz", "Digər"],
    allowCustom: true
  },
  {
    title: "9. Mədə-bağırsaq xəstəlikləri",
    triggerName: "stomach",
    followUps: ["Xora", "Qastrit", "Digər"],
    allowCustom: true
  },
  {
    title: "10. İstifadə edilən dərmanlar",
    triggerName: "medications",
    followUps: ["Antihistamin dərmanlar", "Antiepileptiklər", "İnsulin", "Antihipertenziv", "Hormonal preparatlar", "Ürək dərmanları", "Neyroleptiklər", "Sitostatiklər", "Antikoaqulyantlar", "Digər"],
    allowCustom: true
  },
  {
    title: "11. Hamiləlik",
    triggerName: "pregnancy",
    followUps: ["(həftə)"],
    type: "yes_no_input"
  },
  {
    title: "12. Digər xəstəliklər",
    triggerName: "other_diseases",
    followUps: ["Revmatoid artrit", "Böyrək transplantasiyası", "Radioterapiya", "Sifilis", "Böyrək catışmazlığı", "QICS", "Kimya terapiya", "Hemodializ", "Adrenalin catışmazlığı", "Digər"],
    allowCustom: true
  },
  {
    title: "13. Vərdişlər",
    triggerName: "habits",
    followUps: ["Alkoqol", "Bruksizm", "Siqaret çəkirsinizmi"]
  }
];

const groupContainer = document.getElementById("medicalGroups");
medicalGroups.forEach(group => {
  const groupDiv = document.createElement("div");
  groupDiv.className = "question-group";

  const label = document.createElement("label");
  label.innerText = group.title;
  groupDiv.appendChild(label);

  const yes = document.createElement("label");
  yes.innerHTML = `<input type="radio" name="${group.triggerName}" value="Bəli" required> Bəli`;
  const no = document.createElement("label");
  no.innerHTML = `<input type="radio" name="${group.triggerName}" value="Xeyr"> Xeyr`;

  groupDiv.appendChild(yes);
  groupDiv.appendChild(no);

  const followUpDiv = document.createElement("div");
  followUpDiv.style.marginLeft = "20px";
  followUpDiv.style.display = "none";

  group.followUps.forEach((text, i) => {
    if (text === "(həftə)") {
      const label = document.createElement("label");
      label.innerText = "Neçə həftə?";
      const input = document.createElement("input");
      input.type = "text";
      input.name = `${group.triggerName}_week`;
      followUpDiv.appendChild(label);
      followUpDiv.appendChild(input);
    } else {
      const cbLabel = document.createElement("label");
      cbLabel.innerHTML = `<input type="checkbox" name="${group.triggerName}_sub_${i}" value="${text}"> ${text}`;
      followUpDiv.appendChild(cbLabel);
    }
  });

  if (group.allowCustom) {
    const customInput = document.createElement("input");
    customInput.placeholder = "Digər (əlavə edin)";
    customInput.name = `${group.triggerName}_custom`;
    followUpDiv.appendChild(customInput);
  }

  groupDiv.appendChild(followUpDiv);

  groupDiv.querySelectorAll(`input[name="${group.triggerName}"]`).forEach(radio => {
    radio.addEventListener("change", () => {
      followUpDiv.style.display = radio.value === "Bəli" ? "block" : "none";
    });
  });

  groupContainer.appendChild(groupDiv);
});

document.getElementById("contractForm").onsubmit = async function(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);
  const answers = [];

  for (let [key, val] of formData.entries()) {
    if (val && key.startsWith("group_")) {
      answers.push(`${key}: ${val}`);
    }
  }

  const citizenship = formData.get("citizenship") === "Digər"
    ? document.getElementById("otherCitizenshipInput").value
    : formData.get("citizenship");

  const data = {
    first_name: formData.get("first_name"),
    last_name: formData.get("last_name"),
    middle_name: formData.get("middle_name"),
    birth_date: formData.get("birth_date"),
    phone: formData.get("phone"),
    gender: formData.get("gender"),
    email: formData.get("email"),
    address: formData.get("address"),
    citizenship: citizenship,
    answers: answers
  };

  const file = document.getElementById("photoUpload").files[0];
  if (file) {
    const reader = new FileReader();
    reader.onloadend = async () => {
      data.photo = reader.result;
      await sendForm(data);
    };
    reader.readAsDataURL(file);
  } else {
    await sendForm(data);
  }
};

async function sendForm(data) {
  await fetch("/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data)
  });

  document.getElementById("contractForm").style.display = "none";
  document.getElementById("thanksMessage").style.display = "block";
}
