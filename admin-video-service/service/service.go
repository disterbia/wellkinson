// /admin-video-service/service/service.go
package service

import (
	"admin-video-service/dto"
	"common/model"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type AdminVideoService interface {
	GetLevel1s(id int) ([]dto.VimeoLevel1, error)
	GetLevel2s(id int, projectId string) ([]dto.VimeoLevel2, error)
	SaveVideos(id int, videoIds []string) (string, error)
}

type adminVideoService struct {
	db *gorm.DB
}

func NewAdminVideoService(db *gorm.DB) AdminVideoService {
	return &adminVideoService{db: db}
}

func (service *adminVideoService) GetLevel1s(id int) ([]dto.VimeoLevel1, error) {
	var user model.User
	err := service.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if !user.IsAdmin {
		return nil, errors.New("deny")
	}

	apiURL := "https://api.vimeo.com/users/145953562/projects/14798949/items"

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 생성
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err // 에러 반환
	}

	// Vimeo API 토큰 설정
	req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
	req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // 에러 반환
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // 에러 반환
	}

	var response dto.VimeoResponse
	err = json.Unmarshal(body, &response) // body는 이미 byte slice
	if err != nil {
		return nil, err // 에러 반환
	}

	var vimeoData []dto.VimeoLevel1
	for _, item := range response.Data {
		if item.Type == "folder" {
			// URI에서 프로젝트 ID 추출
			splitUri := strings.Split(item.Folder.Uri, "/")
			projectId := splitUri[len(splitUri)-1]

			// VimeoLevel1 구조체로 데이터 매핑
			vimeoData = append(vimeoData, dto.VimeoLevel1{
				ProjectId: projectId,
				Name:      item.Folder.Name,
			})
		}
	}

	return vimeoData, nil // 결과 반환
}

func (service *adminVideoService) GetLevel2s(id int, projectId string) ([]dto.VimeoLevel2, error) {
	var user model.User
	err := service.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	if !user.IsAdmin {
		return nil, errors.New("deny")
	}

	apiURL := "https://api.vimeo.com/users/145953562/projects/" + projectId + "/videos"

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 생성
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err // 에러 반환
	}

	// Vimeo API 토큰 설정
	req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
	req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

	// 요청 보내기
	resp, err := client.Do(req)

	if err != nil {
		return nil, err // 에러 반환
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // 에러 반환
	}

	var response dto.VimeoResponse2
	err = json.Unmarshal(body, &response) // body는 이미 byte slice
	if err != nil {
		return nil, err // 에러 반환
	}

	var videos []model.Videos
	err = service.db.Where("project_id = ?", projectId).Find(&videos).Error
	if err != nil {
		return nil, errors.New("db error")
	}

	var vimeoData []dto.VimeoLevel2
	for _, item := range response.Data {
		isActive := false
		// URI에서 프로젝트 ID 추출
		splitUri := strings.Split(item.Uri, "/")
		videoId := splitUri[len(splitUri)-1]
		for _, v := range videos {
			if v.VideoId == videoId {
				isActive = true
			}
		}
		vimeoData = append(vimeoData, dto.VimeoLevel2{
			VideoId:  videoId,
			Name:     item.Name,
			IsActive: isActive,
		})

	}

	return vimeoData, nil // 결과 반환
}

func (service *adminVideoService) SaveVideos(id int, videoIds []string) (string, error) {
	var user model.User
	err := service.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return "", errors.New("db error")
	}

	if !user.IsAdmin {
		return "", errors.New("deny")
	}

	// 중복 제거
	seen := make(map[string]bool)
	unique := []string{}

	for _, v := range videoIds {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	// HTTP 클라이언트 생성
	client := &http.Client{}

	var wg sync.WaitGroup
	videosChan := make(chan model.Videos, len(unique))
	proIdChan := make(chan string, len(unique))

	for _, item := range unique {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()

			apiURL := "https://api.vimeo.com/users/145953562/videos/" + item

			// 요청 생성
			req, err := http.NewRequest("GET", apiURL, nil)
			if err != nil {
				log.Println(err)
				return
			}

			// Vimeo API 토큰 설정
			req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
			req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

			// 요청 보내기
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			// 응답 처리
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			var response dto.VimeoResponse3
			err = json.Unmarshal(body, &response) // body는 이미 byte slice
			if err != nil {
				log.Println(err)
				return
			}
			if response.Name == "" {
				log.Println("no id")
				return
			}

			splitUri := strings.Split(response.ParentFolder.Uri, "/")
			projectId := splitUri[len(splitUri)-1]

			videosChan <- model.Videos{
				VideoId:      item,
				Name:         response.Name,
				Duration:     response.Duration,
				ProjectId:    projectId,
				ThumbnailUrl: response.Pictures.BaseLink,
				ProjectName:  response.ParentFolder.Name,
			}
			proIdChan <- projectId
		}(item)
	}

	wg.Wait()
	close(videosChan)
	close(proIdChan)

	if len(videosChan) == 0 {
		return "nothing", nil
	}

	var videos []model.Videos
	var proIds []string
	for video := range videosChan {
		videos = append(videos, video)
	}

	for proId := range proIdChan {
		proIds = append(proIds, proId)
	}

	if err := service.db.Where("project_id IN ?", proIds).Delete(&model.Videos{}).Error; err != nil {
		return "", err
	}

	if err := service.db.Create(&videos).Error; err != nil {
		return "", err

	}

	return "200", nil

}
