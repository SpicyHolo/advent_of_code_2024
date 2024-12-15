import cv2
import numpy as np
import os
import matplotlib.pyplot as plt

def compute_periodicity_metric(image_path):
    # Load the grayscale image
    image = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    
    # Resize for consistency (optional)
    resized_image = cv2.resize(image, (256, 256))
    
    # Compute the DCT of the entire image
    dct_image = cv2.dct(np.float32(resized_image))
    
    # Focus on high-frequency components (e.g., bottom-right quadrant)
    h, w = dct_image.shape
    high_freq_region = dct_image[h//2:, w//2:]
    
    # Compute the periodicity metric: Sum of high-frequency coefficients
    periodicity_metric = np.sum(np.abs(high_freq_region))
    
    return periodicity_metric

def compute_scores_and_plot(image_folder):
    periodicity_scores = []
    image_names = []
    
    # Iterate over all images in the folder
    for image_name in os.listdir(image_folder):
        image_path = os.path.join(image_folder, image_name)
        if image_name.lower().endswith(('.jpg', '.jpeg', '.png', '.bmp', '.tiff')):
            metric = compute_periodicity_metric(image_path)
            periodicity_scores.append(metric)
            image_names.append(image_name)
    
    # Plot the scores
    plt.figure(figsize=(12, 6))
    plt.plot(periodicity_scores, marker='o', linestyle='-', color='b', label='Periodicity Score')
    plt.title('Periodicity Scores Across Images')
    plt.xlabel('Image Index')
    plt.ylabel('Periodicity Score')
    plt.grid(True)
    plt.legend()
    plt.tight_layout()
    
    # Display the plot
    plt.show()
    
    return periodicity_scores, image_names

# Example usage
image_folder = './dst'
scores, names = compute_scores_and_plot(image_folder)