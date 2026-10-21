package postgresql

import (
	"gorm.io/gorm"
)

// Container agrupa y expone todos los repositorios PostgreSQL del servicio.
type Container struct {
	Area              *AreaRepository
	Cargo             *CargoRepository
	Comuna            *ComunaRepository
	Miembro           *MiembroRepository
	MiembroAsignacion *MiembroAsignacionRepository
	Permiso           *PermisoRepository
	Persona           *PersonaRepository
	Provincia         *ProvinciaRepository
	Region            *RegionRepository
	Rol               *RolRepository
	RolPermiso        *RolPermisoRepository
	Sesion            *SesionRepository
	Usuario           *UsuarioRepository
	UsuarioRol        *UsuarioRolRepository
}

// NewContainer inicializa y retorna una nueva instancia contenedora de repositorios PostgreSQL.
func NewContainer(db *gorm.DB) *Container {
	return &Container{
		Area:              NewAreaRepository(db),
		Cargo:             NewCargoRepository(db),
		Comuna:            NewComunaRepository(db),
		Miembro:           NewMiembroRepository(db),
		MiembroAsignacion: NewMiembroAsignacionRepository(db),
		Permiso:           NewPermisoRepository(db),
		Persona:           NewPersonaRepository(db),
		Provincia:         NewProvinciaRepository(db),
		Region:            NewRegionRepository(db),
		Rol:               NewRolRepository(db),
		RolPermiso:        NewRolPermisoRepository(db),
		Sesion:            NewSesionRepository(db),
		Usuario:           NewUsuarioRepository(db),
		UsuarioRol:        NewUsuarioRolRepository(db),
	}
}
