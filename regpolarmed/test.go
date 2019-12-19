package regpolarmed

import "fmt"

//Test func
func Test() {
	fmt.Println("РАЙОНЫ")
	districts, err := GetDistricts()
	if err != nil {
		fmt.Println("Error GetDistricts: ", err)
	}
	for _, d := range districts {
		fmt.Println("Район: " + d.Name + "\t" + "ID: " + d.ID)
	}
	fmt.Println("КЛИНИКИ АПАТИТОВ")
	clinics, err := districts[0].GetClinics()
	if err != nil {
		fmt.Println("Error GetClinics: ", err)
	}
	for _, c := range clinics {
		fmt.Println("ID: ", c.ID)
		fmt.Println("Название: ", c.Name)
		fmt.Println("Адрес: ", c.Address)
	}
	fmt.Println("Специалисты")
	specials, err := clinics[0].GetSpecialists()
	if err != nil {
		fmt.Println("Error GetSpecialists: ", err)
	}
	for _, s := range specials {
		fmt.Println("Специалист: " + s.Name + "\t" + "ID: " + s.ID)
	}
	fmt.Println("Врачи")
	doctors, err := specials[0].GetDoctors()
	if err != nil {
		fmt.Println("Error GetDoctors: ", err)
	}
	for _, d := range doctors {
		fmt.Println("Врач: " + d.Name + "\t" + "ID: " + d.ID)
	}
}
