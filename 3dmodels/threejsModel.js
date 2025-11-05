

import * as THREE from 'three';
import { GLTFLoader } from 'three/examples/jsm/loaders/GLTFLoader';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'; // Импортируем OrbitControls

export default class ThreeModel {
  constructor(container, Model_url, Texture_url, lightMode) {
    this.container = container; // Контейнер для отображения 3D модели
    this.Model_url = Model_url; // Путь к модели
    this.Texture_url = Texture_url; // Путь к текстуре модели
    this.sphereTextureUrl = "../../3dmodels/hdri.png"; // Путь к текстуре для сферы
    this.scene = new THREE.Scene();
    this.camera = new THREE.PerspectiveCamera(75, 1, 0.1, 1000); // Соотношение сторон камеры установлено на 1
    this.renderer = new THREE.WebGLRenderer({ antialias: true });
    this.clock = new THREE.Clock();
    this.model = null;
    this.lightMode = 0; // Текущий режим освещения

    // Инициализация
    console.log(this.sphereTextureUrl);
    this.init();
  }

    // Метод для рандомизации всех параметров
    randomizeAll() {
      this.randomizeLight(); // Рандомизируем свет
      this.randomizeModelScale(); // Рандомизируем масштаб модели
      this.randomizeColor(); // Рандомизируем цвет модели
    }
  
    // Метод для рандомизации интенсивности и цвета света
    randomizeLight() {
  
      const randomIntensity = Math.random() * 2 + 0.5; // Интенсивность от 0.5 до 2
      const randomColor = new THREE.Color(Math.random(), Math.random(), Math.random()); // Случайный цвет
  
      this.updateLightIntensity(randomIntensity);
      this.updateLightColor(randomColor);
    }
  
    // Метод для рандомизации масштаба модели
    randomizeModelScale() {
      const randomScaleX = Math.random() * 3 + 0.5; // Масштаб по оси X от 0.5 до 3
      const randomScaleY = Math.random() * 3 + 0.5; // Масштаб по оси Y от 0.5 до 3
      const randomScaleZ = Math.random() * 3 + 0.5; // Масштаб по оси Z от 0.5 до 3
  
      this.updateModelScaleX(randomScaleX);
      this.updateModelScaleY(randomScaleY);
      this.updateModelScaleZ(randomScaleZ);
    }
  
    // Метод для рандомизации цвета модели
    randomizeColor() {
      if (this.model) {
        const randomColor = new THREE.Color(Math.random(), Math.random(), Math.random()); // Случайный цвет
        this.model.traverse((child) => {
          if (child.isMesh) {
            child.material.color.set(randomColor); // Изменяем цвет материала модели
            child.material.needsUpdate = true;
          }
        });
      }
    }

    // Обновление масштаба модели по оси X
    updateModelScaleX(scaleX) {
      if (this.model) {
        this.model.scale.x = scaleX;
      }
    }
  
    // Обновление масштаба модели по оси Y
    updateModelScaleY(scaleY) {
      if (this.model) {
        this.model.scale.y = scaleY;
      }
    }
  
    // Обновление масштаба модели по оси Z
    updateModelScaleZ(scaleZ) {
      if (this.model) {
        this.model.scale.z = scaleZ;
      }
    }

  updateLightType(type) {
    // // Удаляем текущее освещение
    // this.scene.children = this.scene.children.filter(
    //   (child) => !(child instanceof THREE.Light)
    // );
  
    if (type === 0) {
      // Мягкое освещение
      this.ambientLight = new THREE.AmbientLight(0xffffff, 1);
      this.scene.add(this.ambientLight);
  
      this.directionalLight = new THREE.DirectionalLight(0xffffff, 1);
      this.directionalLight.position.set(1, 1, 1).normalize();
      this.scene.add(this.directionalLight);
    } else if (type === 1) {
      // Точечный источник света
      this.pointLight = new THREE.PointLight(0xffffff, 1);
      this.pointLight.position.set(0, 2, 2);
      this.scene.add(this.pointLight);
    } else if (type === 2) {
      // Цветное освещение
      this.redLight = new THREE.PointLight(0xff0000, 1);
      this.redLight.position.set(-2, 2, 2);
      this.scene.add(this.redLight);
  
      this.greenLight = new THREE.PointLight(0x00ff00, 1);
      this.greenLight.position.set(2, 2, 2);
      this.scene.add(this.greenLight);
  
      this.blueLight = new THREE.PointLight(0x0000ff, 1);
      this.blueLight.position.set(0, -2, 2);
      this.scene.add(this.blueLight);
    }
  }
  
  updateLightIntensity(intensity) {
    // Обновляем интенсивность для всех источников света
    if (this.ambientLight) this.ambientLight.intensity = intensity;
    if (this.directionalLight) this.directionalLight.intensity = intensity;
    if (this.pointLight) this.pointLight.intensity = intensity;
    if (this.redLight) this.redLight.intensity = intensity;
    if (this.greenLight) this.greenLight.intensity = intensity;
    if (this.blueLight) this.blueLight.intensity = intensity;
  }
  
  updateLightColor(color) {
    // Обновляем цвет для всех источников света
    if (this.ambientLight) this.ambientLight.color.set(color);
    if (this.directionalLight) this.directionalLight.color.set(color);
    if (this.pointLight) this.pointLight.color.set(color);
    // Для цветного освещения оставляем текущие цвета
  }

  init() {
    if (!this.container) {
      console.error("Container is not available.");
      return;
    }

    // Устанавливаем размер рендерера на 600x600
    this.renderer.setSize(600, 600);
    this.container.appendChild(this.renderer.domElement); // Добавляем рендерер в контейнер

    // Добавим освещение
    this.addLights();

    // Загрузка модели
    this.loadModel();

    // Загрузка текстуры для сферы
    this.loadSphereTexture();

    // Камера: размещаем ближе
    this.camera.position.z = 1.2;
    this.camera.position.x = 0;
    this.camera.position.y = 0.6;

    // Создаем OrbitControls для управления камерой мышью
    this.controls = new OrbitControls(this.camera, this.renderer.domElement);

    // Создаем кнопку для изменения освещения
    this.createLightToggleButton();

    // Анимация
    this.animate();
  }

  addLights() {
    // Удаляем предыдущее освещение, если есть
    this.scene.children = this.scene.children.filter(
      (child) => !(child instanceof THREE.Light)
    );

    if (this.lightMode === 0) {
      // Мягкое освещение
      const ambientLight = new THREE.AmbientLight(0xffffff, 1);
      this.scene.add(ambientLight);

      const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
      directionalLight.position.set(1, 1, 1).normalize();
      this.scene.add(directionalLight);
    } else if (this.lightMode === 1) {
      // Точечный источник света
      const pointLight = new THREE.PointLight(0xffffff, 1);
      pointLight.position.set(0, 2, 2);
      this.scene.add(pointLight);
    } else if (this.lightMode === 2) {
      // Цветное освещение
      const redLight = new THREE.PointLight(0xff0000, 1, 10);
      redLight.position.set(-2, 2, 2);
      this.scene.add(redLight);

      const greenLight = new THREE.PointLight(0x00ff00, 1, 10);
      greenLight.position.set(2, 2, 2);
      this.scene.add(greenLight);

      const blueLight = new THREE.PointLight(0x0000ff, 1, 10);
      blueLight.position.set(0, -2, 2);
      this.scene.add(blueLight);
    }
  }

    // Добавляем метод для изменения света
    changeLight() {

            // Точечный источник света
            const pointLight = new THREE.PointLight(0xffffff, 1);
            pointLight.position.set(0, 2, 2);
            this.scene.add(pointLight);
  }

  createLightToggleButton() {
    const button = document.createElement("button");
    button.innerText = "Change Light";
    button.style.position = "absolute";
    button.style.top = "10px";
    button.style.left = "10px";
    button.style.padding = "10px";
    button.style.zIndex = "999";

    button.addEventListener("click", () => {
      // Переключаем режим освещения
      this.lightMode = (this.lightMode + 1) % 3;
      this.addLights();
    });

    document.body.appendChild(button);
  }

  loadSphereTexture() {
    const textureLoader = new THREE.TextureLoader();

    // Загружаем текстуру для сферы
    textureLoader.load(
      this.sphereTextureUrl,
      (texture) => {
        // Создаем большую сферу с текстурой
        const geometry = new THREE.SphereGeometry(500, 60, 40);
        geometry.scale(-1, 1, 1); // Отображаем сферу изнутри
        const material = new THREE.MeshBasicMaterial({
          map: texture,
          side: THREE.DoubleSide, // Чтобы текстура была видна изнутри
        });
        const sphere = new THREE.Mesh(geometry, material);
        this.scene.add(sphere);
      },
      undefined,
      (error) => {
        console.error("Ошибка загрузки текстуры сферы:", error);
      }
    );
  }

  loadModel() {
    const loader = new GLTFLoader();

      // Проверка на undefined или пустую строку
  if (!this.Model_url || typeof this.Model_url !== 'string') {
    console.error('Неверный URL модели:', this.Model_url);
    return;
  }

    // Загружаем модель, используя переданный путь
    loader.load(
      this.Model_url,
      (gltf) => {
        this.model = gltf.scene;

        // Загрузка текстуры модели
        if (this.Texture_url) {
          const textureLoader = new THREE.TextureLoader();
          const texture = textureLoader.load(this.Texture_url, () => {
            // Применяем текстуру ко всем материалам модели
            this.model.traverse((child) => {
              if (child.isMesh) {
                child.material.map = texture;
                child.material.needsUpdate = true;
              }
            });
          });
        }

        this.scene.add(this.model);
      },
      undefined,
      (error) => {
        console.error(error);
      }
    );
  }

  animate() {
    requestAnimationFrame(() => this.animate());

    // Обновляем контроллеры (для управления мышью)
    this.controls.update();

    this.renderer.render(this.scene, this.camera);
  }

  resize() {
    // Обновляем соотношение сторон камеры для фиксированного размера 600x600
    this.camera.aspect = 1; // Соотношение сторон теперь 1 (600x600)
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(600, 600);
  }
}