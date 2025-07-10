// usecase/raza_usecase.go
package usecase

import (
	"context" // Necesario si tus métodos de repositorio usan context.Context

	"github.com/BastianCarrasco/backend-go/model"
	"github.com/BastianCarrasco/backend-go/repository"
)

// RazaUseCase define la interfaz para las operaciones de negocio de Raza.
// Esto es útil para la inyección de dependencias y pruebas unitarias.
type RazaUseCase interface {
	GetRazaByID(ctx context.Context, id string) (*model.Raza, error)
	GetAllRaces(ctx context.Context) ([]model.Raza, error)
	// Puedes añadir más métodos como CreateRaza, UpdateRaza, DeleteRaza
}

// razaUseCase implementa la interfaz RazaUseCase.
type razaUseCase struct {
	razaRepo repository.RazaRepository
	// Puedes añadir otros repositorios o servicios aquí si son necesarios
}

// NewRazaUseCase crea una nueva instancia de RazaUseCase.
func NewRazaUseCase(repo repository.RazaRepository) RazaUseCase {
	return &razaUseCase{
		razaRepo: repo,
	}
}

// GetRazaByID implementa la lógica para obtener una raza por su ID.
func (uc *razaUseCase) GetRazaByID(ctx context.Context, id string) (*model.Raza, error) {
	// Aquí podrías añadir lógica de negocio adicional antes o después de llamar al repositorio.
	// Por ejemplo, validaciones, transformaciones de datos, etc.
	raza, err := uc.razaRepo.FindRazaByID(id)
	if err != nil {
		return nil, err
	}
	// Podrías verificar si raza es nil aquí si FindRazaByID retorna (nil, nil) cuando no lo encuentra
	// if raza == nil {
	// 	return nil, errors.New("raza no encontrada") // Ejemplo de un error más específico
	// }
	return raza, nil
}

// GetAllRaces implementa la lógica para obtener todas las razas.
func (uc *razaUseCase) GetAllRaces(ctx context.Context) ([]model.Raza, error) {
	// Aquí podrías añadir lógica de negocio adicional, como filtrado, paginación, etc.
	razas, err := uc.razaRepo.GetAllRaces()
	if err != nil {
		return nil, err
	}
	return razas, nil
}