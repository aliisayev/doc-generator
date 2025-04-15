document.addEventListener("DOMContentLoaded", function () {
    console.log("DOMContentLoaded event started");
    const medicalGroups = [
      {
        title: "Ürək-damar xəstəlikləri",
        triggerName: "heart_diseases",
        followUps: [
          "Revmakardit",
          "Miokard infarktı",
          "Oynaqların kəskin revmatizmi",
          "Böyrək xəstəlikləri",
          "Stenokardiya",
          "Ürək çatışmazlığı",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Arterial təzyiq",
        triggerName: "blood_pressure",
        followUps: ["Hipertenziya", "Hipotenziya"],
        allowCustom: false
      },
      {
        title: "Qan sisteminin xəstəlikləri",
        triggerName: "blood_diseases",
        followUps: [
          "Hemorragik diatez",
          "Anemiya",
          "Laxtalanma sisteminin anomaliyası",
          "Leykemiya",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "QBB xəstəlikləri",
        triggerName: "ent_diseases",
        followUps: [
          "Frontal Sinusit",
          "Maksillar Sinusit",
          "Tonzil Hipertrofiyası",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Ağciyər xəstəlikləri",
        triggerName: "lung_diseases",
        followUps: [
          "Astma",
          "Keçirilmiş tuberkulyoz",
          "Aktiv tuberkulyoz",
          "Bronxit",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Psixo-nevroloji xəstəliklər",
        triggerName: "neuro_diseases",
        followUps: [
          "Epilepsiya",
          "Nevralgiyalar",
          "Paraliç",
          "Psixiatrik müalicə",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Qaraciyər xəstəlikləri",
        triggerName: "liver_diseases",
        followUps: [
          "Hepatit A",
          "Hepatit B",
          "Hepatit C",
          "Qaraciyər çatışmazlığı",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Endokrin xəstəlikləri",
        triggerName: "endocrine_diseases",
        followUps: [
          "Diabet",
          "Hipotireoz",
          "Paratireoz",
          "Hipertireoz",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Mədə-bağırsaq xəstəlikləri",
        triggerName: "gi_diseases",
        followUps: ["Xora", "Qastrit", "Digər"],
        allowCustom: true
      },
      {
        title: "İstifadə edilən dərmanlar",
        triggerName: "medications",
        followUps: [
          "Antihistamin dərmanlar",
          "Antiepileptiklər",
          "İnsulin",
          "Antihipertenziv",
          "Hormonal preparatlar",
          "Ürək dərmanları",
          "Neyroleptiklər",
          "Sitostatiklər",
          "Antikoaqulyantlar",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Hamiləlik",
        triggerName: "pregnancy",
        followUps: ["(həftə)"],
        allowCustom: false
      },
      {
        title: "Digər xəstəliklər",
        triggerName: "other_diseases",
        followUps: [
          "Revmatoid artrit",
          "Böyrək transplantasiyası",
          "Radioterapiya",
          "Sifilis",
          "Böyrək çatışmazlığı",
          "QİÇS",
          "Kimya terapiya",
          "Hemodializ",
          "Adrenalin çatışmazlığı",
          "Digər"
        ],
        allowCustom: true
      },
      {
        title: "Vərdişlər",
        triggerName: "habits",
        followUps: ["Alkoqol", "Bruksizm", "Siqaret çəkirsinizmi"],
        allowCustom: false
      }
    ];
    console.log("medicalGroups array defined");
  
    const groupContainer = document.getElementById("medicalGroups");
  
    medicalGroups.forEach(group => {
      console.log("Starting iteration for group:", group.title);
      const groupDiv = document.createElement("div");
      groupDiv.className = "question-group";
      groupDiv.setAttribute("data-trigger", group.triggerName);
  
      const label = document.createElement("label");
      label.innerText = group.title;
      groupDiv.appendChild(label);
  
      const yes = document.createElement("label");
      yes.innerHTML = `Bəli <input type="radio" name="${group.triggerName}" value="Bəli" required>`;
      const no = document.createElement("label");
      no.innerHTML = `Xeyr <input type="radio" name="${group.triggerName}" value="Xeyr">`;
  
      groupDiv.appendChild(yes);
      groupDiv.appendChild(no);
  
      const followUpDiv = document.createElement("div");
      followUpDiv.style.marginTop = "10px";
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
          const subDiv = document.createElement("div");
          subDiv.className = "subquestion";
  
          const subLabel = document.createElement("div");
          subLabel.innerText = text;
          subDiv.appendChild(subLabel);
  
          const bəliRadio = document.createElement("label");
          bəliRadio.innerHTML = `Bəli <input type="radio" name="${group.triggerName}_sub_${i}" value="Bəli">`;
  
          const xeyrRadio = document.createElement("label");
          xeyrRadio.innerHTML = `Xeyr <input type="radio" name="${group.triggerName}_sub_${i}" value="Xeyr">`;
  
          subDiv.appendChild(bəliRadio);
          subDiv.appendChild(xeyrRadio);
          followUpDiv.appendChild(subDiv);
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
  
    function togglePregnancyBlock() {
      console.log("togglePregnancyBlock function called");
      const genderSelect = document.querySelector('select[name="gender"]');
      const gender = genderSelect.value;
      const pregnancyBlock = document.querySelector('[data-trigger="pregnancy"]');
  
      if (pregnancyBlock) {
        const inputs = pregnancyBlock.querySelectorAll('input[name="pregnancy"]');
        if (gender === "Qadın") {
          pregnancyBlock.style.display = "block";
          inputs.forEach(input => input.setAttribute("required", "required"));
        } else {
          pregnancyBlock.style.display = "none";
          inputs.forEach(input => input.removeAttribute("required"));
        }
      }
      console.log("togglePregnancyBlock function finished");
    }
  
    document.querySelector('select[name="gender"]').addEventListener("change", togglePregnancyBlock);
    togglePregnancyBlock();
  
    document.getElementById("citizenshipSelect").addEventListener("change", function () {
      console.log("citizenshipSelect event listener triggered");
      const isOther = this.value === "Digər";
      document.getElementById("otherCitizenshipField").style.display = isOther ? "block" : "none";
    });
  
    document.getElementById("photoUpload").addEventListener("change", function () {
      console.log("photoUpload event listener triggered");
      const file = this.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = () => {
          const img = new Image();
          img.src = reader.result;
          img.style.maxWidth = "100%";
          document.getElementById("photoPreview").innerHTML = "";
          document.getElementById("photoPreview").appendChild(img);
        };
        reader.readAsDataURL(file);
      }
      console.log("photoUpload event listener finished");
    });
  
    document.getElementById("contractForm").onsubmit = async function (e) {
      console.log("contractForm submit event triggered");
      e.preventDefault();
  
      const form = e.target;
      const firstInvalid = form.querySelector(":invalid");
      if (firstInvalid) {
        firstInvalid.scrollIntoView({ behavior: "smooth", block: "center" });
        firstInvalid.focus();
        return;
      }
      console.log("Form validation passed");
  
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
      console.log("sendForm function called with data:", data);
      await fetch("/submit", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
      });
      console.log("fetch request completed");
  
      document.getElementById("contractForm").style.display = "none";
      document.getElementById("thanksMessage").style.display = "block";
    }
  });