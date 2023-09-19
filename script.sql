-- Creación de la base de datos 'clinic'
CREATE DATABASE IF NOT EXISTS `my_db` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `my_db`;

CREATE TABLE IF NOT EXISTS `odontologo` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Identificador del odontologo en el sistema',
  `apellido` VARCHAR(100) NOT NULL COMMENT 'Apellido del odontologo',
  `nombre` VARCHAR(100) NOT NULL COMMENT 'Nombre del odontologo',
  `matricula` VARCHAR(100) NOT NULL COMMENT 'Número de licencia del odontologo',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARACTER SET = utf8mb3;

CREATE TABLE IF NOT EXISTS `paciente` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Identificador paciente en el sistema',
  `nombre` VARCHAR(100) NOT NULL COMMENT 'Nombre del paciente',
  `apellido` VARCHAR(100) NOT NULL COMMENT 'Apellido del paciente',
  `domicilio` VARCHAR(100) NULL DEFAULT NULL COMMENT 'Dirección del paciente',
  `dni` INT NOT NULL COMMENT 'Identificación del paciente',
  `fecha_alta` DATE NOT NULL COMMENT 'Fecha de alta del paciente',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARACTER SET = utf8mb3;

CREATE TABLE IF NOT EXISTS `turno` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Identificador del turno',
  `id_odontologo` INT NULL DEFAULT NULL COMMENT 'Identificador del odontólogo',
  `id_paciente` INT NOT NULL COMMENT 'Identificador del paciente',
  `fecha_hora` DATETIME NULL DEFAULT NULL COMMENT 'Fecha y hora del turno',
  `descripcion` VARCHAR(300) NULL DEFAULT NULL COMMENT 'Descripcion del turno',
  PRIMARY KEY (`id`),
  INDEX `turno_FK` (`id_odontologo` ASC) VISIBLE,
  INDEX `turno_FK_1` (`id_paciente` ASC) VISIBLE,
  CONSTRAINT `turno_FK`
    FOREIGN KEY (`id`)
    REFERENCES `odontologo` (`id`),
  CONSTRAINT `turno_FK_1`
    FOREIGN KEY (`id`)
    REFERENCES `paciente` (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARACTER SET = utf8mb3;

-- Inserciones en la tabla 'odontologo'
INSERT INTO `odontologo` (`apellido`, `nombre`, `matricula`)
VALUES
('Pérez', 'Juan', '12345'),
('Gómez', 'María', '67890'),
('López', 'Carlos', '54321');

-- Inserciones en la tabla 'paciente'
INSERT INTO `paciente` (`nombre`, `apellido`, `domicilio`, `dni`, `fecha_alta`)
VALUES
('Ana', 'Martínez', 'Calle 123', 12345678, '2023-09-18'),
('Pedro', 'González', 'Avenida XYZ', 87654321, '2023-09-19'),
('Laura', 'Díaz', 'Calle ABC', 98765432, '2023-09-20');


-- Inserciones en la tabla 'turno'
INSERT INTO `turno` (`id_odontologo`, `id_paciente`, `fecha_hora`, `descripcion`)
VALUES
(1, 1, '2023-09-22 10:00:00.000', 'Limpieza dental'),
(2, 2, '2023-09-23 15:30:00.000', 'Extracción de muelas'),
(3, 3, '2023-09-24 09:15:00.000', 'Consulta de rutina');