package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"core-service/internal/enums"
	modelsRedis "core-service/internal/models/redis"
	"core-service/internal/services"
	pb "core-service/proto"
)

// AccesoHandler implementa el servidor gRPC para el servicio AccesoService.
type AccesoHandler struct {
	pb.UnimplementedAccesoServiceServer
	service *services.AccesoService
}

// NewAccesoHandler inicializa y retorna una nueva instancia de AccesoHandler.
func NewAccesoHandler(service *services.AccesoService) *AccesoHandler {
	return &AccesoHandler{
		service: service,
	}
}

// IniciarSesion atiende la solicitud de autenticación mediante gRPC.
func (h *AccesoHandler) IniciarSesion(ctx context.Context, req *pb.IniciarSesionRequest) (*pb.IniciarSesionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "la solicitud no puede ser nula")
	}

	input := services.LoginInput{
		Run:               req.Run,
		Pasaporte:         req.Pasaporte,
		CorreoElectronico: req.CorreoElectronico,
		Contrasena:        req.Contrasena,
	}

	output, err := h.service.Login(ctx, input)
	if err != nil {
		if err.Error() == "credenciales inválidas" {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.IniciarSesionResponse{
		Token:     output.Token,
		TokenType: output.TokenType,
		ExpiresIn: output.ExpiresIn,
		Sesion:    toProtoSesionData(output.Sesion),
	}, nil
}

// CerrarSesion revoca la sesión activa en el sistema.
func (h *AccesoHandler) CerrarSesion(ctx context.Context, req *pb.CerrarSesionRequest) (*pb.CerrarSesionResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "el token es obligatorio")
	}

	if err := h.service.Logout(ctx, req.Token); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.CerrarSesionResponse{
		Exito:   true,
		Mensaje: "sesión cerrada exitosamente",
	}, nil
}

// CerrarTodasLasSesiones revoca todas las sesiones activas del usuario en Redis.
func (h *AccesoHandler) CerrarTodasLasSesiones(ctx context.Context, req *pb.CerrarTodasLasSesionesRequest) (*pb.CerrarTodasLasSesionesResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "el token es obligatorio")
	}

	if err := h.service.LogoutAll(ctx, req.Token); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.CerrarTodasLasSesionesResponse{
		Exito:   true,
		Mensaje: "todas las sesiones activas han sido revocadas exitosamente",
	}, nil
}

// ObtenerPerfil retorna los datos de la sesión activa recuperados directamente desde la caché.
func (h *AccesoHandler) ObtenerPerfil(ctx context.Context, req *pb.ObtenerPerfilRequest) (*pb.ObtenerPerfilResponse, error) {
	if req == nil || req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "el token es obligatorio")
	}

	sesion, err := h.service.Me(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.ObtenerPerfilResponse{
		Sesion: toProtoSesionData(sesion),
	}, nil
}

// ActualizarContrasena actualiza la clave secreta del usuario autenticado.
func (h *AccesoHandler) ActualizarContrasena(ctx context.Context, req *pb.ActualizarContrasenaRequest) (*pb.ActualizarContrasenaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "la solicitud no puede ser nula")
	}

	input := services.UpdatePasswordInput{
		ContrasenaActual: req.ContrasenaActual,
		NuevaContrasena:  req.NuevaContrasena,
	}

	if err := h.service.UpdatePassword(ctx, req.Token, input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.ActualizarContrasenaResponse{
		Exito:   true,
		Mensaje: "contraseña actualizada exitosamente",
	}, nil
}

// RegistrarUsuario realiza el autoregistro de una nueva persona y usuario.
func (h *AccesoHandler) RegistrarUsuario(ctx context.Context, req *pb.RegistrarUsuarioRequest) (*pb.RegistrarUsuarioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "la solicitud no puede ser nula")
	}

	var comunaID *uuid.UUID
	if req.ComunaId != nil && *req.ComunaId != "" {
		parsed, err := uuid.Parse(*req.ComunaId)
		if err == nil {
			comunaID = &parsed
		}
	}

	var fechaNacimiento *time.Time
	if req.FechaNacimiento != nil && *req.FechaNacimiento != "" {
		parsed, err := time.Parse("2006-01-02", *req.FechaNacimiento)
		if err == nil {
			fechaNacimiento = &parsed
		}
	}

	var genero *enums.Genero
	if req.Genero != nil && *req.Genero != "" {
		g := enums.Genero(*req.Genero)
		genero = &g
	}

	input := services.RegisterInput{
		Run:               req.Run,
		Pasaporte:         req.Pasaporte,
		Nombre:            req.Nombre,
		ApellidoPaterno:   req.ApellidoPaterno,
		ApellidoMaterno:   req.ApellidoMaterno,
		Genero:            genero,
		FechaNacimiento:   fechaNacimiento,
		ComunaID:          comunaID,
		Direccion:         req.Direccion,
		TelefonoFijo:      req.TelefonoFijo,
		TelefonoMovil:     req.TelefonoMovil,
		CorreoElectronico: req.CorreoElectronico,
		Contrasena:        req.Contrasena,
	}

	usuario, err := h.service.Register(ctx, input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RegistrarUsuarioResponse{
		UsuarioId: usuario.ID.String(),
		PersonaId: usuario.PersonaID.String(),
		Mensaje:   "usuario registrado exitosamente",
	}, nil
}

// SolicitarRecuperacionContrasena procesa la solicitud de código de recuperación.
func (h *AccesoHandler) SolicitarRecuperacionContrasena(ctx context.Context, req *pb.SolicitarRecuperacionContrasenaRequest) (*pb.SolicitarRecuperacionContrasenaResponse, error) {
	if req == nil || req.CorreoElectronico == "" {
		return nil, status.Error(codes.InvalidArgument, "el correo electrónico es obligatorio")
	}

	input := services.ForgotPasswordInput{
		CorreoElectronico: req.CorreoElectronico,
	}

	_ = h.service.ForgotPassword(ctx, input)

	return &pb.SolicitarRecuperacionContrasenaResponse{
		Exito:   true,
		Mensaje: "si el correo se encuentra registrado, se enviará un enlace de recuperación",
	}, nil
}

// RestablecerContrasena valida el código de recuperación y asigna la nueva contraseña.
func (h *AccesoHandler) RestablecerContrasena(ctx context.Context, req *pb.RestablecerContrasenaRequest) (*pb.RestablecerContrasenaResponse, error) {
	if req == nil || req.Token == "" || req.NuevaContrasena == "" {
		return nil, status.Error(codes.InvalidArgument, "el código y la nueva contraseña son obligatorios")
	}

	input := services.ResetPasswordInput{
		Token:           req.Token,
		NuevaContrasena: req.NuevaContrasena,
	}

	if err := h.service.ResetPassword(ctx, input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.RestablecerContrasenaResponse{
		Exito:   true,
		Mensaje: "contraseña restablecida exitosamente",
	}, nil
}

// toProtoSesionData convierte el modelo de sesión de Redis a su representación Protobuf.
func toProtoSesionData(sesion *modelsRedis.Sesion) *pb.SesionData {
	if sesion == nil {
		return nil
	}

	p := sesion.Persona
	protoPersona := &pb.SesionPersona{
		Id:                p.ID.String(),
		Run:               p.Run,
		Pasaporte:         p.Pasaporte,
		Nombre:            p.Nombre,
		ApellidoPaterno:   p.ApellidoPaterno,
		ApellidoMaterno:   p.ApellidoMaterno,
		Direccion:         p.Direccion,
		TelefonoFijo:      p.TelefonoFijo,
		TelefonoMovil:     p.TelefonoMovil,
		CorreoElectronico: p.CorreoElectronico,
	}

	if p.Genero != nil {
		g := p.Genero.String()
		protoPersona.Genero = &g
	}

	if p.FechaNacimiento != nil {
		fn := p.FechaNacimiento.Format("2006-01-02")
		protoPersona.FechaNacimiento = &fn
	}

	if p.Comuna != nil {
		protoPersona.Comuna = &pb.SesionComuna{
			Id:          p.Comuna.ID.String(),
			Nombre:      p.Comuna.Nombre,
			ProvinciaId: p.Comuna.ProvinciaID.String(),
		}
	}

	var protoMiembro *pb.SesionMiembro
	if sesion.Miembro != nil {
		protoMiembro = &pb.SesionMiembro{
			Id:                sesion.Miembro.ID.String(),
			TelefonoFijo:      sesion.Miembro.TelefonoFijo,
			TelefonoMovil:     sesion.Miembro.TelefonoMovil,
			CorreoElectronico: sesion.Miembro.CorreoElectronico,
		}
	}

	protoUsuario := &pb.SesionUsuario{
		Id: sesion.Usuario.ID.String(),
	}

	var protoAplicaciones []*pb.SesionAplicacion
	for _, app := range sesion.Aplicaciones {
		var protoRoles []*pb.SesionRol
		for _, rol := range app.Roles {
			var protoPermisos []*pb.SesionPermiso
			for _, permiso := range rol.Permisos {
				protoPermisos = append(protoPermisos, &pb.SesionPermiso{
					Id:     permiso.ID.String(),
					Nombre: permiso.Nombre,
					Slug:   permiso.Slug,
				})
			}
			protoRoles = append(protoRoles, &pb.SesionRol{
				Id:       rol.ID.String(),
				Nombre:   rol.Nombre,
				Slug:     rol.Slug,
				Permisos: protoPermisos,
			})
		}
		protoAplicaciones = append(protoAplicaciones, &pb.SesionAplicacion{
			AplicacionSlug: app.AplicacionSlug.String(),
			Roles:          protoRoles,
		})
	}

	return &pb.SesionData{
		Persona:      protoPersona,
		Miembro:      protoMiembro,
		Usuario:      protoUsuario,
		Aplicaciones: protoAplicaciones,
	}
}
