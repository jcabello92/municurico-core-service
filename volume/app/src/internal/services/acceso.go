package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"core-service/internal/enums"
	"core-service/internal/interceptors"
	modelsPostgres "core-service/internal/models/postgresql"
	modelsRedis "core-service/internal/models/redis"
	reposPostgres "core-service/internal/repositories/postgresql"
	reposRedis "core-service/internal/repositories/redis"
	"core-service/internal/security"
)

// LoginInput contiene los datos de entrada requeridos para la autenticación.
type LoginInput struct {
	Run               *string `json:"run,omitempty"`
	Pasaporte         *string `json:"pasaporte,omitempty"`
	CorreoElectronico *string `json:"correo_electronico,omitempty"`
	Contrasena        string  `json:"contrasena"`
}

// LoginOutput contiene los datos del token JWT generado y la sesión activa resultante.
type LoginOutput struct {
	Token     string              `json:"token"`
	TokenType string              `json:"token_type"`
	ExpiresIn int64               `json:"expires_in"`
	Sesion    *modelsRedis.Sesion `json:"sesion"`
}

// LogoutInput contiene el token de la sesión a cerrar.
type LogoutInput struct {
	Token string `json:"token"`
}

// UpdatePasswordInput contiene las contraseñas requeridas para el cambio seguro de credenciales.
type UpdatePasswordInput struct {
	ContrasenaActual string `json:"contrasena_actual"`
	NuevaContrasena  string `json:"nueva_contrasena"`
}

// RegisterInput contiene los datos personales y de cuenta requeridos para registrar un nuevo usuario.
type RegisterInput struct {
	Run               *string       `json:"run,omitempty"`
	Pasaporte         *string       `json:"pasaporte,omitempty"`
	Nombre            string        `json:"nombre"`
	ApellidoPaterno   *string       `json:"apellido_paterno,omitempty"`
	ApellidoMaterno   *string       `json:"apellido_materno,omitempty"`
	Genero            *enums.Genero `json:"genero,omitempty"`
	FechaNacimiento   *time.Time    `json:"fecha_nacimiento,omitempty"`
	ComunaID          *uuid.UUID    `json:"comuna_id,omitempty"`
	Direccion         *string       `json:"direccion,omitempty"`
	TelefonoFijo      *string       `json:"telefono_fijo,omitempty"`
	TelefonoMovil     *string       `json:"telefono_movil,omitempty"`
	CorreoElectronico string        `json:"correo_electronico"`
	Contrasena        string        `json:"contrasena"`
}

// ForgotPasswordInput contiene el correo electrónico para solicitar la recuperación de contraseña.
type ForgotPasswordInput struct {
	CorreoElectronico string `json:"correo_electronico"`
}

// ResetPasswordInput contiene el token de recuperación y la nueva contraseña a establecer.
type ResetPasswordInput struct {
	Token           string `json:"token"`
	NuevaContrasena string `json:"nueva_contrasena"`
}

// AccesoService coordina la lógica de negocio para autenticación, autorización y administración de sesiones.
type AccesoService struct {
	repos                 *reposPostgres.Container
	redisSesionRepo       *reposRedis.SesionRepository
	redisRecuperacionRepo *reposRedis.RecuperacionRepository
	hasher                *security.PasswordHasher
	encryptor             *security.Encryptor
	bindexer              *security.Bindexer
	jwt                   *security.JWT
	sessionTTL            time.Duration
}

// NewAccesoService inicializa y retorna una nueva instancia de AccesoService con todas sus dependencias.
func NewAccesoService(
	repos *reposPostgres.Container,
	redisSesionRepo *reposRedis.SesionRepository,
	redisRecuperacionRepo *reposRedis.RecuperacionRepository,
	hasher *security.PasswordHasher,
	encryptor *security.Encryptor,
	bindexer *security.Bindexer,
	jwt *security.JWT,
	sessionTTL time.Duration,
) *AccesoService {
	return &AccesoService{
		repos:                 repos,
		redisSesionRepo:       redisSesionRepo,
		redisRecuperacionRepo: redisRecuperacionRepo,
		hasher:                hasher,
		encryptor:             encryptor,
		bindexer:              bindexer,
		jwt:                   jwt,
		sessionTTL:            sessionTTL,
	}
}

// Login autentica a un usuario mediante sus credenciales y genera una nueva sesión en Redis junto con su token JWT.
func (s *AccesoService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	var count int
	if input.Run != nil && strings.TrimSpace(*input.Run) != "" {
		count++
	}
	if input.Pasaporte != nil && strings.TrimSpace(*input.Pasaporte) != "" {
		count++
	}
	if input.CorreoElectronico != nil && strings.TrimSpace(*input.CorreoElectronico) != "" {
		count++
	}

	if count == 0 {
		return nil, errors.New("debe ingresar al menos un identificador (run, pasaporte o correo electrónico)")
	}
	if count > 1 {
		return nil, errors.New("debe ingresar únicamente un identificador (run, pasaporte o correo electrónico)")
	}
	if strings.TrimSpace(input.Contrasena) == "" {
		return nil, errors.New("la contraseña no puede estar vacía")
	}

	var persona *modelsPostgres.Persona
	var err error

	if input.Run != nil && strings.TrimSpace(*input.Run) != "" {
		bindex, bErr := s.bindexer.GenerateBindex(*input.Run)
		if bErr != nil {
			return nil, fmt.Errorf("error al generar blind index para run: %w", bErr)
		}
		persona, err = s.repos.Persona.ReadByRunBindex(ctx, bindex)
	} else if input.Pasaporte != nil && strings.TrimSpace(*input.Pasaporte) != "" {
		bindex, bErr := s.bindexer.GenerateBindex(*input.Pasaporte)
		if bErr != nil {
			return nil, fmt.Errorf("error al generar blind index para pasaporte: %w", bErr)
		}
		persona, err = s.repos.Persona.ReadByPasaporteBindex(ctx, bindex)
	} else if input.CorreoElectronico != nil && strings.TrimSpace(*input.CorreoElectronico) != "" {
		bindex, bErr := s.bindexer.GenerateBindex(*input.CorreoElectronico)
		if bErr != nil {
			return nil, fmt.Errorf("error al generar blind index para correo electrónico: %w", bErr)
		}
		persona, err = s.repos.Persona.ReadByCorreoElectronicoBindex(ctx, bindex)
	}

	if err != nil || persona == nil {
		return nil, errors.New("credenciales inválidas")
	}

	usuario, err := s.repos.Usuario.ReadByPersonaID(ctx, persona.ID.String())
	if err != nil || usuario == nil {
		return nil, errors.New("credenciales inválidas")
	}

	if usuario.Estado != enums.EstadoUsuarioActivo {
		return nil, fmt.Errorf("el usuario no está activo (estado: %s)", usuario.Estado)
	}

	if usuario.Contrasena == nil || !s.hasher.CheckPassword(*usuario.Contrasena, input.Contrasena) {
		return nil, errors.New("credenciales inválidas")
	}

	sesion, err := s.buildRedisSesion(ctx, persona, usuario)
	if err != nil {
		return nil, fmt.Errorf("error al construir sesión: %w", err)
	}

	now := time.Now().UTC()
	timestamp := now.Unix()

	if err := s.redisSesionRepo.Save(ctx, sesion, timestamp, s.sessionTTL); err != nil {
		return nil, fmt.Errorf("error al almacenar sesión en redis: %w", err)
	}

	// Registrar sesión en PostgreSQL para trazabilidad
	audit := interceptors.GetAuditFromContext(ctx)
	ip := "127.0.0.1"
	if audit.AutorIP != nil && *audit.AutorIP != "" {
		ip = *audit.AutorIP
	}
	ua := "desconocido"
	if audit.UserAgent != nil && *audit.UserAgent != "" {
		ua = *audit.UserAgent
	}

	dbSesion := &modelsPostgres.Sesion{
		UsuarioID:       usuario.ID,
		UsuarioIP:       ip,
		UserAgent:       ua,
		FechaInicio:     now,
		FechaExpiracion: now.Add(s.sessionTTL),
		Estado:          enums.EstadoSesionActiva,
	}
	_, _ = s.repos.Sesion.Create(ctx, dbSesion)

	sesionKey := sesion.BuildKey(timestamp)
	token, err := s.jwt.Generate(sesionKey, s.sessionTTL)
	if err != nil {
		return nil, fmt.Errorf("error al generar token jwt: %w", err)
	}

	return &LoginOutput{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(s.sessionTTL.Seconds()),
		Sesion:    sesion,
	}, nil
}

// Logout revoca la sesión activa eliminando la clave de Redis.
func (s *AccesoService) Logout(ctx context.Context, tokenString string) error {
	if strings.TrimSpace(tokenString) == "" {
		return errors.New("el token no puede estar vacío")
	}

	claims, err := s.jwt.Validate(tokenString)
	if err != nil {
		return fmt.Errorf("token inválido: %w", err)
	}

	if err := s.redisSesionRepo.Delete(ctx, claims.SesionKey); err != nil {
		return fmt.Errorf("error al eliminar sesión de redis: %w", err)
	}

	return nil
}

// LogoutAll revoca de forma masiva todas las sesiones activas del usuario en Redis.
func (s *AccesoService) LogoutAll(ctx context.Context, tokenString string) error {
	if strings.TrimSpace(tokenString) == "" {
		return errors.New("el token no puede estar vacío")
	}

	claims, err := s.jwt.Validate(tokenString)
	if err != nil {
		return fmt.Errorf("token inválido: %w", err)
	}

	sesion, err := s.redisSesionRepo.Get(ctx, claims.SesionKey)
	if err != nil {
		return fmt.Errorf("sesión no encontrada o expirada: %w", err)
	}

	if err := s.redisSesionRepo.DeleteAllByUsuarioID(ctx, sesion.Usuario.ID); err != nil {
		return fmt.Errorf("error al revocar sesiones del usuario: %w", err)
	}

	return nil
}

// Me recupera los datos de la sesión activa directamente desde Redis utilizando el token provisto.
func (s *AccesoService) Me(ctx context.Context, tokenString string) (*modelsRedis.Sesion, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("el token no puede estar vacío")
	}

	claims, err := s.jwt.Validate(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	sesion, err := s.redisSesionRepo.Get(ctx, claims.SesionKey)
	if err != nil {
		return nil, errors.New("la sesión ha expirado o no existe")
	}

	return sesion, nil
}

// UpdatePassword actualiza la contraseña del usuario y revoca sus sesiones previas en Redis.
func (s *AccesoService) UpdatePassword(ctx context.Context, tokenString string, input UpdatePasswordInput) error {
	if strings.TrimSpace(tokenString) == "" {
		return errors.New("el token no puede estar vacío")
	}
	if strings.TrimSpace(input.ContrasenaActual) == "" || strings.TrimSpace(input.NuevaContrasena) == "" {
		return errors.New("la contraseña actual y la nueva contraseña son obligatorias")
	}
	if len(input.NuevaContrasena) < 8 {
		return errors.New("la nueva contraseña debe tener al menos 8 caracteres")
	}

	claims, err := s.jwt.Validate(tokenString)
	if err != nil {
		return fmt.Errorf("token inválido: %w", err)
	}

	sesion, err := s.redisSesionRepo.Get(ctx, claims.SesionKey)
	if err != nil {
		return errors.New("sesión expirada o no encontrada")
	}

	usuario, err := s.repos.Usuario.Read(ctx, sesion.Usuario.ID.String())
	if err != nil {
		return fmt.Errorf("error al recuperar usuario: %w", err)
	}

	if usuario.Contrasena == nil || !s.hasher.CheckPassword(*usuario.Contrasena, input.ContrasenaActual) {
		return errors.New("la contraseña actual no coincide")
	}

	newHash, err := s.hasher.HashPassword(input.NuevaContrasena)
	if err != nil {
		return fmt.Errorf("error al generar hash de la nueva contraseña: %w", err)
	}

	if err := s.repos.Usuario.UpdateContrasena(ctx, usuario.ID.String(), newHash); err != nil {
		return fmt.Errorf("error al actualizar contraseña en base de datos: %w", err)
	}

	_ = s.redisSesionRepo.DeleteAllByUsuarioID(ctx, usuario.ID)
	now := time.Now().UTC().Unix()
	_ = s.redisSesionRepo.Save(ctx, sesion, now, s.sessionTTL)

	return nil
}

// Register registra una nueva entidad Persona y su respectivo Usuario en la base de datos.
func (s *AccesoService) Register(ctx context.Context, input RegisterInput) (*modelsPostgres.Usuario, error) {
	if strings.TrimSpace(input.Nombre) == "" {
		return nil, errors.New("el nombre es obligatorio")
	}
	if strings.TrimSpace(input.CorreoElectronico) == "" {
		return nil, errors.New("el correo electrónico es obligatorio")
	}
	if strings.TrimSpace(input.Contrasena) == "" {
		return nil, errors.New("la contraseña es obligatoria")
	}
	if len(input.Contrasena) < 8 {
		return nil, errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	correoBindex, err := s.bindexer.GenerateBindex(input.CorreoElectronico)
	if err != nil {
		return nil, fmt.Errorf("error al generar blind index para correo electrónico: %w", err)
	}

	existente, _ := s.repos.Persona.ReadByCorreoElectronicoBindex(ctx, correoBindex)
	if existente != nil {
		return nil, errors.New("ya existe un usuario registrado con ese correo electrónico")
	}

	correoCifrado, err := s.encryptor.Encrypt(input.CorreoElectronico)
	if err != nil {
		return nil, fmt.Errorf("error al cifrar correo electrónico: %w", err)
	}

	persona := &modelsPostgres.Persona{
		Nombre:                  strings.TrimSpace(input.Nombre),
		ApellidoPaterno:         input.ApellidoPaterno,
		ApellidoMaterno:         input.ApellidoMaterno,
		Genero:                  input.Genero,
		FechaNacimiento:         input.FechaNacimiento,
		ComunaID:                input.ComunaID,
		CorreoElectronico:       &correoCifrado,
		CorreoElectronicoBindex: &correoBindex,
	}

	if input.Run != nil && strings.TrimSpace(*input.Run) != "" {
		runBindex, err := s.bindexer.GenerateBindex(*input.Run)
		if err != nil {
			return nil, fmt.Errorf("error al generar blind index para run: %w", err)
		}
		runCifrado, err := s.encryptor.Encrypt(*input.Run)
		if err != nil {
			return nil, fmt.Errorf("error al cifrar run: %w", err)
		}
		persona.Run = &runCifrado
		persona.RunBindex = &runBindex
	}

	if input.Pasaporte != nil && strings.TrimSpace(*input.Pasaporte) != "" {
		pasaporteBindex, err := s.bindexer.GenerateBindex(*input.Pasaporte)
		if err != nil {
			return nil, fmt.Errorf("error al generar blind index para pasaporte: %w", err)
		}
		pasaporteCifrado, err := s.encryptor.Encrypt(*input.Pasaporte)
		if err != nil {
			return nil, fmt.Errorf("error al cifrar pasaporte: %w", err)
		}
		persona.Pasaporte = &pasaporteCifrado
		persona.PasaporteBindex = &pasaporteBindex
	}

	if input.Direccion != nil && strings.TrimSpace(*input.Direccion) != "" {
		dirCifrada, err := s.encryptor.Encrypt(*input.Direccion)
		if err == nil {
			persona.Direccion = &dirCifrada
		}
	}
	if input.TelefonoFijo != nil && strings.TrimSpace(*input.TelefonoFijo) != "" {
		tfCifrado, err := s.encryptor.Encrypt(*input.TelefonoFijo)
		if err == nil {
			persona.TelefonoFijo = &tfCifrado
		}
	}
	if input.TelefonoMovil != nil && strings.TrimSpace(*input.TelefonoMovil) != "" {
		tmCifrado, err := s.encryptor.Encrypt(*input.TelefonoMovil)
		if err == nil {
			persona.TelefonoMovil = &tmCifrado
		}
	}

	personaIDStr, err := s.repos.Persona.Create(ctx, persona)
	if err != nil {
		return nil, fmt.Errorf("error al crear persona: %w", err)
	}

	personaUUID, err := uuid.Parse(personaIDStr)
	if err != nil {
		return nil, fmt.Errorf("error al procesar id de persona creada: %w", err)
	}

	hashedPassword, err := s.hasher.HashPassword(input.Contrasena)
	if err != nil {
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	usuario := &modelsPostgres.Usuario{
		PersonaID:  personaUUID,
		Contrasena: &hashedPassword,
		Estado:     enums.EstadoUsuarioActivo,
		EsSistema:  false,
	}

	usuarioIDStr, err := s.repos.Usuario.Create(ctx, usuario)
	if err != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", err)
	}

	return s.repos.Usuario.Read(ctx, usuarioIDStr)
}

// ForgotPassword genera y almacena un token temporal en Redis para el restablecimiento de contraseña.
func (s *AccesoService) ForgotPassword(ctx context.Context, input ForgotPasswordInput) error {
	if strings.TrimSpace(input.CorreoElectronico) == "" {
		return errors.New("el correo electrónico no puede estar vacío")
	}

	correoBindex, err := s.bindexer.GenerateBindex(input.CorreoElectronico)
	if err != nil {
		return fmt.Errorf("error al generar blind index: %w", err)
	}

	persona, err := s.repos.Persona.ReadByCorreoElectronicoBindex(ctx, correoBindex)
	if err != nil || persona == nil {
		// Retornar sin error para evitar enumeración de correos
		return nil
	}

	usuario, err := s.repos.Usuario.ReadByPersonaID(ctx, persona.ID.String())
	if err != nil || usuario == nil {
		return nil
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("error al generar token aleatorio: %w", err)
	}
	tokenString := hex.EncodeToString(tokenBytes)

	recuperacion := &modelsRedis.Recuperacion{
		UsuarioID: usuario.ID,
		CreatedAt: time.Now().UTC(),
	}
	ttl := 15 * time.Minute

	if err := s.redisRecuperacionRepo.Save(ctx, tokenString, recuperacion, ttl); err != nil {
		return fmt.Errorf("error al guardar token de recuperación en redis: %w", err)
	}

	return nil
}

// ResetPassword valida el token de recuperación y actualiza la contraseña del usuario.
func (s *AccesoService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	if strings.TrimSpace(input.Token) == "" {
		return errors.New("el token de recuperación no puede estar vacío")
	}
	if strings.TrimSpace(input.NuevaContrasena) == "" {
		return errors.New("la nueva contraseña no puede estar vacía")
	}
	if len(input.NuevaContrasena) < 8 {
		return errors.New("la nueva contraseña debe tener al menos 8 caracteres")
	}

	recuperacion, err := s.redisRecuperacionRepo.GetByToken(ctx, input.Token)
	if err != nil || recuperacion == nil || recuperacion.UsuarioID == uuid.Nil {
		return errors.New("el código de recuperación es inválido o ha expirado")
	}

	newHash, err := s.hasher.HashPassword(input.NuevaContrasena)
	if err != nil {
		return fmt.Errorf("error al hashear nueva contraseña: %w", err)
	}

	if err := s.repos.Usuario.UpdateContrasena(ctx, recuperacion.UsuarioID.String(), newHash); err != nil {
		return fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	_ = s.redisRecuperacionRepo.DeleteByToken(ctx, input.Token)
	_ = s.redisSesionRepo.DeleteAllByUsuarioID(ctx, recuperacion.UsuarioID)

	return nil
}

// buildRedisSesion construye la estructura de sesión Redis con datos descifrados y roles agrupados por aplicación.
func (s *AccesoService) buildRedisSesion(
	ctx context.Context,
	persona *modelsPostgres.Persona,
	usuario *modelsPostgres.Usuario,
) (*modelsRedis.Sesion, error) {
	sesionPersona := modelsRedis.SesionPersona{
		ID:                persona.ID,
		Run:               s.decryptStringPtr(persona.Run),
		Pasaporte:         s.decryptStringPtr(persona.Pasaporte),
		Nombre:            persona.Nombre,
		ApellidoPaterno:   persona.ApellidoPaterno,
		ApellidoMaterno:   persona.ApellidoMaterno,
		Genero:            persona.Genero,
		FechaNacimiento:   persona.FechaNacimiento,
		Direccion:         s.decryptStringPtr(persona.Direccion),
		TelefonoFijo:      s.decryptStringPtr(persona.TelefonoFijo),
		TelefonoMovil:     s.decryptStringPtr(persona.TelefonoMovil),
		CorreoElectronico: s.decryptStringPtr(persona.CorreoElectronico),
	}

	if persona.ComunaID != nil {
		comuna, err := s.repos.Comuna.Read(ctx, persona.ComunaID.String())
		if err == nil && comuna != nil {
			sesionPersona.Comuna = &modelsRedis.SesionComuna{
				ID:          comuna.ID,
				Nombre:      comuna.Nombre,
				ProvinciaID: comuna.ProvinciaID,
			}
		}
	}

	var sesionMiembro *modelsRedis.SesionMiembro
	miembro, err := s.repos.Miembro.ReadByPersonaID(ctx, persona.ID.String())
	if err == nil && miembro != nil && miembro.Estado == enums.EstadoMiembroActivo {
		sesionMiembro = &modelsRedis.SesionMiembro{
			ID:                miembro.ID,
			TelefonoFijo:      s.decryptStringPtr(miembro.TelefonoFijo),
			TelefonoMovil:     s.decryptStringPtr(miembro.TelefonoMovil),
			CorreoElectronico: s.decryptStringPtr(miembro.CorreoElectronico),
		}
	}

	sesionUsuario := modelsRedis.SesionUsuario{
		ID: usuario.ID,
	}

	usuarioRoles, err := s.repos.UsuarioRol.ListByUsuarioID(ctx, usuario.ID.String())
	appsMap := make(map[enums.AplicacionSlug][]modelsRedis.SesionRol)

	if err == nil {
		for _, ur := range usuarioRoles {
			rol, rErr := s.repos.Rol.Read(ctx, ur.RolID.String())
			if rErr != nil || rol == nil {
				continue
			}

			rolPermisos, _ := s.repos.RolPermiso.ListByRolID(ctx, ur.RolID.String())
			var permisos []modelsRedis.SesionPermiso

			for _, rp := range rolPermisos {
				permiso, pErr := s.repos.Permiso.Read(ctx, rp.PermisoID.String())
				if pErr == nil && permiso != nil {
					permisos = append(permisos, modelsRedis.SesionPermiso{
						ID:     permiso.ID,
						Nombre: permiso.Nombre,
						Slug:   permiso.Slug,
					})
				}
			}

			appsMap[rol.AplicacionSlug] = append(appsMap[rol.AplicacionSlug], modelsRedis.SesionRol{
				ID:       rol.ID,
				Nombre:   rol.Nombre,
				Slug:     rol.Slug,
				Permisos: permisos,
			})
		}
	}

	var aplicaciones []modelsRedis.SesionAplicacion
	for appSlug, roles := range appsMap {
		aplicaciones = append(aplicaciones, modelsRedis.SesionAplicacion{
			AplicacionSlug: appSlug,
			Roles:          roles,
		})
	}

	return &modelsRedis.Sesion{
		Persona:      sesionPersona,
		Miembro:      sesionMiembro,
		Usuario:      sesionUsuario,
		Aplicaciones: aplicaciones,
	}, nil
}

// decryptStringPtr descifra un puntero a string cifrado retornando un puntero al texto plano.
func (s *AccesoService) decryptStringPtr(encrypted *string) *string {
	if encrypted == nil || *encrypted == "" {
		return nil
	}
	decrypted, err := s.encryptor.Decrypt(*encrypted)
	if err != nil {
		return nil
	}
	return &decrypted
}
