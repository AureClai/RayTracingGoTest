package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	geom "./geometry"
	"./view"
)

// WIDTH unexported
const WIDTH int = 400

// HEIGHT unexported
const HEIGHT int = 400

// SAMPLES unexported
const SAMPLES int = 100

func drand48() float64 {
	return rand.Float64()
}

func writeLines(lines []string, file *os.File) error {
	for i := 0; i < len(lines); i++ {
		_, err := fmt.Fprintln(file, lines[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func color(r *geom.Ray, world geom.Hitable, lightShape geom.Hitable, depth int) *geom.Vec3 {
	var hrec = geom.HitRecord{}

	if world.Hit(r, 0.001, math.MaxFloat64, &hrec) {
		srec := geom.ScatterRecord{}
		var emitted = new(geom.Vec3)
		var scattered = new(geom.Ray)
		pdfVal := 0.0
		*emitted = *hrec.MatPtr.Emitted(r, &hrec, hrec.U, hrec.V, hrec.P)
		if depth < 50 && hrec.MatPtr.Scatter(r, &hrec, &srec) {
			if srec.IsSpecular {
				return srec.Attenuation.Times(color(srec.SpecularRay, world, lightShape, depth+1))
			} else {
				pLight := geom.NewHitablePdf(lightShape, hrec.P)
				p := geom.NewMixturePdf(pLight, srec.PdfPtr)
				*scattered = *geom.NewRayWithTime(hrec.P, p.Generate(), r.Time())
				pdfVal = p.Value(scattered.Direction())
				//fmt.Println("Other")
				return ((color(scattered, world, lightShape, depth+1).Times(srec.Attenuation.TimesScalar(hrec.MatPtr.ScatteringPdf(r, &hrec, scattered)))).Plus(emitted)).ByScalar(pdfVal)
			}
		}
		return emitted
	}
	return geom.NewVec3(0.0, 0.0, 0.0)
}

func MakecornellBoxObjects() *geom.HitableList {
	list := make([]geom.Hitable, 50)
	red := geom.Lambertian{geom.NewConstantTexture(geom.NewVec3(0.65, 0.05, 0.05))}
	white := geom.Lambertian{geom.NewConstantTexture(geom.NewVec3(0.73, 0.73, 0.73))}
	green := geom.Lambertian{geom.NewConstantTexture(geom.NewVec3(0.12, 0.45, 0.15))}
	light := geom.DiffuseLight{Emit: geom.NewConstantTexture(geom.NewVec3(15, 15, 15))}
	list[0] = geom.NewFlipNormals(geom.NewYZRect(0, 555, 0, 555, 555, green))
	list[1] = geom.NewYZRect(0, 555, 0, 555, 0, red)
	list[2] = geom.NewFlipNormals(geom.NewXZRect(213, 343, 227, 332, 554, light))
	list[3] = geom.NewFlipNormals(geom.NewXZRect(0, 555, 0, 555, 555, white))
	list[4] = geom.NewXZRect(0, 555, 0, 555, 0, white)
	list[5] = geom.NewFlipNormals(geom.NewXYRect(0, 555, 0, 555, 555, white))
	// boxes
	glass := geom.Dielectric{1.5}
	//list[6] = geom.NewTranslate(geom.NewRotateY(geom.NewBox(geom.NewVec3(0, 0, 0), geom.NewVec3(165, 165, 165), white), -18), geom.NewVec3(130, 0, 65))
	list[6] = geom.NewSphere(geom.NewVec3(190, 90, 190), 90, glass)
	//aluminium := geom.Metal{Albedo: geom.NewVec3(0.8, 0.85, 0.88), Fuzz: 0.0}
	list[7] = geom.NewTranslate(geom.NewRotateY(geom.NewBox(geom.NewVec3(0, 0, 0), geom.NewVec3(165, 330, 165), white), 15), geom.NewVec3(265, 0, 295))
	return geom.NewHitableList(&list, 8)
}

type Scene struct {
	Objects  *geom.HitableList
	Camera   *view.Camera
	Settings *SceneSettings
}

func cornellBox() *Scene {
	settings := &SceneSettings{
		Width:   WIDTH,
		Height:  HEIGHT,
		Samples: SAMPLES,
	}
	var lookFrom = geom.NewVec3(278, 278, -800)
	var lookAt = geom.NewVec3(278, 278, 0)
	distToFocus := 10.0
	aperture := 0.0
	vfov := 40.0
	aspect := float64(settings.Width) / float64(settings.Height)
	var cam = view.NewCamera(lookFrom, lookAt, geom.NewVec3(0, 1, 0), vfov, aspect, aperture, distToFocus, 0.0, 1.0)
	scene := MakecornellBoxObjects()
	return &Scene{
		Objects:  scene,
		Camera:   cam,
		Settings: settings,
	}
}

type SceneSettings struct {
	Width   int
	Height  int
	Samples int
}

func deNan(v *geom.Vec3) *geom.Vec3 {
	x := v.X()
	y := v.Y()
	z := v.Z()
	if x < 0 || math.IsNaN(x) {
		x = 0.0
	}
	if y < 0 || math.IsNaN(y) {
		y = 0.0
	}
	if y < 0 || math.IsNaN(y) {
		y = 0.0
	}
	return geom.NewVec3(x, y, z)
}

func main() {
	// time management
	start := time.Now()
	percentage := 0.0

	f, err := os.Create("test3.ppm")
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	//
	scene := cornellBox()
	//
	var lines = make([]string, scene.Settings.Width*scene.Settings.Height+3)
	// Header of Picture
	lines[0] = "P3"
	lines[1] = fmt.Sprintf("%v %v", scene.Settings.Width, scene.Settings.Height)
	lines[2] = fmt.Sprintf("%v", 255)

	// Lines
	for j := scene.Settings.Height - 1; j >= 0; j-- {
		// Columns
		for i := 0; i < scene.Settings.Width; i++ {
			percentage = 100.0 * float64(scene.Settings.Width*(scene.Settings.Height-1-j)+i) / float64(scene.Settings.Width*scene.Settings.Height)
			fmt.Printf("\r%5.2f %%", percentage)
			var col = geom.NewVec3(0, 0, 0)
			lightShape := geom.NewXZRect(213, 343, 227, 332, 554, geom.NewNoMaterial())
			glassSphere := geom.NewSphere(geom.NewVec3(190, 90, 190), 90, geom.NewNoMaterial())
			list := make([]geom.Hitable, 2)
			list[0] = lightShape
			list[1] = glassSphere
			hList := geom.NewHitableList(&list, 2)
			// Samples
			for s := 0; s < SAMPLES; s++ {
				u := (float64(i) + rand.Float64()) / float64(scene.Settings.Width)
				v := (float64(j) + rand.Float64()) / float64(scene.Settings.Height-1)
				var r = scene.Camera.GetRay(u, v)
				col = col.Plus(deNan(color(r, scene.Objects, hList, 0)))
			}
			col = col.ByScalar(float64(SAMPLES))
			col = geom.NewVec3(math.Sqrt(col.R()), math.Sqrt(col.G()), math.Sqrt(col.B()))
			ir := int(255.99 * col.R())
			ig := int(255.99 * col.G())
			ib := int(255.99 * col.B())
			lines[3+scene.Settings.Width*(scene.Settings.Height-1-j)+i] = fmt.Sprintf("%v %v %v", ir, ig, ib)
		}
	}

	// Actually writying the file
	err = writeLines(lines, f)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Closing the file
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Time management and display of end
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Calculations made in %v\n", elapsed)
	fmt.Println("file written successfully")
}
